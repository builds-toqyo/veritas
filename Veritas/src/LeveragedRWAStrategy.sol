// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "./VeritasInvoiceToken.sol";

/**
 * @title IMantleLendingProtocol
 * @notice Interface for Mantle-based lending protocols (INIT Capital, Lendle)
 * @dev Allows borrowing stablecoins against mETH collateral
 */
interface IMantleLendingProtocol {
    function supply(address asset, uint256 amount) external;
    function borrow(address asset, uint256 amount) external;
    function repay(address asset, uint256 amount) external;
    function withdraw(address asset, uint256 amount) external;
    
    function getAccountLiquidity(address account) external view returns (
        uint256 collateralValue,
        uint256 borrowValue,
        uint256 healthFactor
    );
    
    function getCollateralFactor(address asset) external view returns (uint256);
    function getBorrowRate(address asset) external view returns (uint256);
}

// ============================================================================
// LEVERAGE STRATEGY: Borrow Against mETH for RWA Exposure
// ============================================================================

/**
 * @title LeveragedRWAStrategy
 * @notice Composable strategy: mETH collateral → borrow USDC → deploy to RWA
 * @dev Demonstrates DeFi composability and lending integration
 * 
 * FLOW:
 * 1. Vault deposits mETH as collateral to Mantle lending protocol
 * 2. Borrow USDC against mETH (at safe LTV, e.g., 60%)
 * 3. Deploy borrowed USDC to RWA Strategy (AIT tokens)
 * 4. Earn RWA yield (8% APY target)
 * 5. Use yield to repay USDC borrow and compound
 * 
 * BENEFIT: Leveraged exposure to RWA without selling mETH (maintains LST yield)
 */
contract LeveragedRWAStrategy is AccessControl {
    
    bytes32 public constant VAULT_ROLE = keccak256("VAULT_ROLE");
    bytes32 public constant KEEPER_ROLE = keccak256("KEEPER_ROLE");
    
    IERC20 public immutable mETH;
    IERC20 public immutable usdc;
    IMantleLendingProtocol public immutable lendingProtocol;
    VeritasInvoiceToken public immutable ait;
    
    // Leverage parameters
    uint256 public targetLTV = 6000; // 60% LTV (basis points)
    uint256 public maxLTV = 7000;    // 70% max (safety buffer)
    uint256 public constant MAX_BPS = 10000;
    
    // Position tracking
    uint256 public totalCollateral;  // mETH supplied
    uint256 public totalBorrowed;    // USDC borrowed
    uint256 public totalAITHoldings; // AIT tokens held
    
    // Risk metrics
    uint256 public currentHealthFactor; // Updated on each action
    uint256 public minHealthFactor = 13000; // 1.3x minimum
    
    event CollateralSupplied(uint256 mETHAmount, uint256 timestamp);
    event StablecoinBorrowed(uint256 usdcAmount, uint256 newLTV);
    event RWADeployed(uint256 usdcAmount, uint256 aitReceived);
    event YieldHarvested(uint256 usdcYield);
    event LeverageReduced(uint256 repayAmount, string reason);
    event HealthFactorUpdated(uint256 oldHF, uint256 newHF);
    
    constructor(
        address _mETH,
        address _usdc,
        address _lendingProtocol,
        address _ait
    ) {
        mETH = IERC20(_mETH);
        usdc = IERC20(_usdc);
        lendingProtocol = IMantleLendingProtocol(_lendingProtocol);
        ait = VeritasInvoiceToken(_ait);
        
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(VAULT_ROLE, msg.sender);
    }
    
    /**
     * @notice Supply mETH as collateral to lending protocol
     * @param amount mETH to supply
     */
    function supplyCollateral(uint256 amount) 
        external 
        onlyRole(VAULT_ROLE) 
    {
        require(mETH.transferFrom(msg.sender, address(this), amount), "Transfer failed");
        
        // Supply to lending protocol
        mETH.approve(address(lendingProtocol), amount);
        lendingProtocol.supply(address(mETH), amount);
        
        totalCollateral += amount;
        
        _updateHealthFactor();
        
        emit CollateralSupplied(amount, block.timestamp);
    }
    
    /**
     * @notice Borrow USDC against mETH collateral
     * @param amount USDC to borrow
     */
    function borrowStablecoin(uint256 amount) 
        external 
        onlyRole(KEEPER_ROLE) 
    {
        // Check LTV
        uint256 newBorrowed = totalBorrowed + amount;
        uint256 newLTV = (newBorrowed * MAX_BPS) / totalCollateral;
        
        require(newLTV <= targetLTV, "Exceeds target LTV");
        
        // Borrow from protocol
        lendingProtocol.borrow(address(usdc), amount);
        totalBorrowed += amount;
        
        _updateHealthFactor();
        
        emit StablecoinBorrowed(amount, newLTV);
    }
    
    /**
     * @notice Deploy borrowed USDC to RWA (buy AIT tokens)
     * @param usdcAmount Amount to deploy
     */
    function deployToRWA(uint256 usdcAmount) 
        external 
        onlyRole(KEEPER_ROLE) 
    {
        require(usdc.balanceOf(address(this)) >= usdcAmount, "Insufficient USDC");
        
        // Calculate AIT to receive based on NAV
        uint256 aitToReceive = (usdcAmount * 1e6) / ait.navPerToken();
        
        // In production: would execute swap on Mantle DEX or direct mint
        // For simulation: assume direct mint at NAV
        usdc.approve(address(ait), usdcAmount);
        ait.mint(address(this), aitToReceive);
        
        totalAITHoldings += aitToReceive;
        
        emit RWADeployed(usdcAmount, aitToReceive);
    }
    
    /**
     * @notice Harvest yield from AIT holdings
     * @dev Yield comes from invoice payments distributed to AIT holders
     */
    function harvestRWAYield() 
        external 
        onlyRole(KEEPER_ROLE) 
        returns (uint256 yieldAmount) 
    {
        // Get current NAV
        uint256 currentValue = (totalAITHoldings * ait.navPerToken()) / 1e6;
        uint256 initialValue = (totalBorrowed * MAX_BPS) / targetLTV; // Approx initial
        
        if (currentValue > initialValue) {
            yieldAmount = currentValue - initialValue;
            
            // Redeem portion of AIT for USDC
            uint256 aitToBurn = (yieldAmount * 1e6) / ait.navPerToken();
            ait.burn(address(this), aitToBurn);
            totalAITHoldings -= aitToBurn;
            
            // Now have USDC yield
            emit YieldHarvested(yieldAmount);
        }
        
        return yieldAmount;
    }
    
    /**
     * @notice Repay USDC debt using harvested yield
     * @param amount USDC to repay
     */
    function repayDebt(uint256 amount) 
        external 
        onlyRole(KEEPER_ROLE) 
    {
        require(usdc.balanceOf(address(this)) >= amount, "Insufficient USDC");
        
        usdc.approve(address(lendingProtocol), amount);
        lendingProtocol.repay(address(usdc), amount);
        
        totalBorrowed -= amount;
        
        _updateHealthFactor();
    }
    
    /**
     * @notice Emergency deleverage if health factor drops
     * @dev Sells AIT for USDC and repays debt
     */
    function emergencyDeleverage(uint256 aitToSell) 
        external 
        onlyRole(KEEPER_ROLE) 
    {
        require(currentHealthFactor < minHealthFactor, "Health factor OK");
        
        // Sell AIT for USDC
        uint256 usdcReceived = (aitToSell * ait.navPerToken()) / 1e6;
        ait.burn(address(this), aitToSell);
        totalAITHoldings -= aitToSell;
        
        // Repay debt
        usdc.approve(address(lendingProtocol), usdcReceived);
        lendingProtocol.repay(address(usdc), usdcReceived);
        totalBorrowed -= usdcReceived;
        
        _updateHealthFactor();
        
        emit LeverageReduced(usdcReceived, "Emergency deleverage");
    }
    
    /**
     * @notice Update health factor from lending protocol
     */
    function _updateHealthFactor() internal {
        (uint256 collateralValue, uint256 borrowValue, uint256 hf) = 
            lendingProtocol.getAccountLiquidity(address(this));
        
        uint256 oldHF = currentHealthFactor;
        currentHealthFactor = hf;
        
        emit HealthFactorUpdated(oldHF, hf);
    }
    
    /**
     * @notice Get current leverage metrics
     */
    function getLeverageMetrics() external view returns (
        uint256 ltv,
        uint256 healthFactor,
        uint256 aitValue,
        uint256 netExposure
    ) {
        ltv = totalCollateral > 0 ? (totalBorrowed * MAX_BPS) / totalCollateral : 0;
        healthFactor = currentHealthFactor;
        aitValue = (totalAITHoldings * ait.navPerToken()) / 1e6;
        netExposure = aitValue > totalBorrowed ? aitValue - totalBorrowed : 0;
    }
}
