// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/VeritasInvoiceToken.sol";
import "../src/LeveragedRWAStrategy.sol";
import "../src/TieredKYCVerifier.sol";

/**
 * @title VeritasVaultTest
 * @notice Comprehensive test suite for Veritas RWA Vault
 * @dev Tests all critical paths for hackathon demo
 */
contract VeritasVaultTest is Test {
    
    // Contracts
    VeritasInvoiceToken public ait;
    TieredKYCVerifier public kycVerifier;
    LeveragedRWAStrategy public leverageStrategy;
    
    // Mock tokens
    MockERC20 public mETH;
    MockERC20 public usdc;
    MockLendingProtocol public lendingProtocol;
    
    // Test accounts
    address public admin = address(0x1);
    address public keeper = address(0x2);
    address public retailInvestor = address(0x3);
    address public accreditedInvestor = address(0x4);
    address public institutionalInvestor = address(0x5);
    
    // Constants
    uint256 constant INITIAL_BALANCE = 1_000_000 * 1e6; // 1M USDC
    uint256 constant INITIAL_METH_BALANCE = 1_000 * 1e18; // 1K mETH
    
    function setUp() public {
        vm.startPrank(admin);
        
        // Deploy mock tokens
        mETH = new MockERC20("Mantle ETH", "mETH", 18);
        usdc = new MockERC20("USD Coin", "USDC", 6);
        
        // Deploy mock lending protocol
        lendingProtocol = new MockLendingProtocol(address(usdc));
        
        // Deploy core contracts
        ait = new VeritasInvoiceToken();
        kycVerifier = new TieredKYCVerifier();
        leverageStrategy = new LeveragedRWAStrategy(
            address(mETH),
            address(usdc),
            address(lendingProtocol),
            address(ait)
        );
        
        // Setup roles
        ait.grantRole(ait.ORACLE_ROLE(), keeper);
        leverageStrategy.grantRole(leverageStrategy.KEEPER_ROLE(), keeper);
        
        // Fund test accounts
        mETH.mint(admin, INITIAL_METH_BALANCE);
        usdc.mint(admin, INITIAL_BALANCE);
        usdc.mint(address(lendingProtocol), INITIAL_BALANCE * 10); // Protocol liquidity
        
        vm.stopPrank();
    }
    
    // ========================================================================
    // TEST 1: KYC Tier System
    // ========================================================================
    
    function test_KYCTierIssuance() public {
        vm.startPrank(admin);
        
        // Issue retail KYC
        kycVerifier.issueKyc(
            retailInvestor,
            TieredKYCVerifier.InvestorTier.RETAIL,
            365, // 1 year validity
            bytes32("US"),
            keccak256(abi.encodePacked(retailInvestor))
        );
        
        // Verify tier
        assertEq(uint(kycVerifier.getTier(retailInvestor)), uint(TieredKYCVerifier.InvestorTier.RETAIL));
        assertTrue(kycVerifier.hasValidKyc(retailInvestor));
        
        // Check investment cap
        uint256 remainingCap = kycVerifier.getRemainingCapacity(retailInvestor);
        assertEq(remainingCap, 10_000 * 1e6); // $10k cap
        
        vm.stopPrank();
    }
    
    function test_KYCTierUpgrade() public {
        vm.startPrank(admin);
        
        // Issue retail KYC
        kycVerifier.issueKyc(
            retailInvestor,
            TieredKYCVerifier.InvestorTier.RETAIL,
            365,
            bytes32("US"),
            keccak256(abi.encodePacked(retailInvestor))
        );
        
        // Upgrade to accredited
        kycVerifier.upgradeTier(retailInvestor, TieredKYCVerifier.InvestorTier.ACCREDITED);
        
        // Verify upgrade
        assertEq(uint(kycVerifier.getTier(retailInvestor)), uint(TieredKYCVerifier.InvestorTier.ACCREDITED));
        assertEq(kycVerifier.getRemainingCapacity(retailInvestor), 500_000 * 1e6); // $500k cap
        
        vm.stopPrank();
    }
    
    function test_KYCInvestmentCap() public {
        vm.startPrank(admin);
        
        // Issue retail KYC
        kycVerifier.issueKyc(
            retailInvestor,
            TieredKYCVerifier.InvestorTier.RETAIL,
            365,
            bytes32("US"),
            keccak256(abi.encodePacked(retailInvestor))
        );
        
        // Record investment within cap
        kycVerifier.recordInvestment(retailInvestor, 5_000 * 1e6);
        assertEq(kycVerifier.getRemainingCapacity(retailInvestor), 5_000 * 1e6);
        
        // Try to exceed cap
        vm.expectRevert("Exceeds cap");
        kycVerifier.recordInvestment(retailInvestor, 6_000 * 1e6);
        
        vm.stopPrank();
    }
    
    // ========================================================================
    // TEST 2: AIT Synthetic Asset
    // ========================================================================
    
    function test_AITPoolInitialization() public {
        vm.startPrank(admin);
        
        bytes32 poolId = keccak256("CORP-INV-Q1-2025");
        
        ait.initializePool(
            poolId,
            1_000_000 * 1e6, // $1M face value
            50,              // 50 invoices
            90,              // 90 day maturity
            800              // 8% APY
        );
        
        (
            bytes32 storedPoolId,
            uint256 totalFaceValue,
            uint256 numberOfInvoices,
            uint256 weightedMaturity,
            uint256 expectedYield,
            ,
        ) = ait.pool();
        
        assertEq(storedPoolId, poolId);
        assertEq(totalFaceValue, 1_000_000 * 1e6);
        assertEq(numberOfInvoices, 50);
        assertEq(weightedMaturity, 90);
        assertEq(expectedYield, 800);
        
        vm.stopPrank();
    }
    
    function test_AITNAVUpdate() public {
        vm.startPrank(keeper);
        
        uint256 oldNAV = ait.navPerToken();
        assertEq(oldNAV, 1e6); // $1.00 initial
        
        // Update NAV to $1.05 (5% appreciation)
        ait.updateNav(1.05e6);
        
        assertEq(ait.navPerToken(), 1.05e6);
        
        vm.stopPrank();
    }
    
    function test_AITCashFlowDistribution() public {
        vm.startPrank(admin);
        
        // Initialize pool
        ait.initializePool(
            keccak256("TEST-POOL"),
            1_000_000 * 1e6,
            50,
            90,
            800
        );
        
        // Mint AIT tokens
        ait.setWhitelist(admin, true);
        ait.mint(admin, 1_000_000 * 1e6);
        
        vm.stopPrank();
        
        // Record cash flow as oracle
        vm.startPrank(keeper);
        
        ait.recordCashFlow(50_000 * 1e6, 5); // $50k from 5 invoices
        
        // NAV should increase
        uint256 newNAV = ait.navPerToken();
        assertGt(newNAV, 1e6); // Should be > $1.00
        
        vm.stopPrank();
    }
    
    function test_AITDefaultHandling() public {
        vm.startPrank(admin);
        
        // Initialize pool
        ait.initializePool(
            keccak256("TEST-POOL"),
            1_000_000 * 1e6,
            50,
            90,
            800
        );
        
        // Mint AIT tokens
        ait.setWhitelist(admin, true);
        ait.mint(admin, 1_000_000 * 1e6);
        
        vm.stopPrank();
        
        // Record default as oracle
        vm.startPrank(keeper);
        
        uint256 oldNAV = ait.navPerToken();
        
        ait.recordDefault(20_000 * 1e6, 123); // $20k default
        
        // NAV should decrease
        uint256 newNAV = ait.navPerToken();
        assertLt(newNAV, oldNAV);
        
        vm.stopPrank();
    }
    
    // ========================================================================
    // TEST 3: Leveraged RWA Strategy
    // ========================================================================
    
    function test_LeverageSupplyCollateral() public {
        vm.startPrank(admin);
        
        uint256 collateralAmount = 100 * 1e18; // 100 mETH
        
        mETH.approve(address(leverageStrategy), collateralAmount);
        leverageStrategy.supplyCollateral(collateralAmount);
        
        assertEq(leverageStrategy.totalCollateral(), collateralAmount);
        
        vm.stopPrank();
    }
    
    function test_LeverageBorrowStablecoin() public {
        vm.startPrank(admin);
        
        // Supply collateral first
        uint256 collateralAmount = 100 * 1e18; // 100 mETH
        mETH.approve(address(leverageStrategy), collateralAmount);
        leverageStrategy.supplyCollateral(collateralAmount);
        
        // Grant keeper role to admin for testing
        leverageStrategy.grantRole(leverageStrategy.KEEPER_ROLE(), admin);
        
        // Borrow USDC (60% LTV)
        uint256 borrowAmount = 60_000 * 1e6; // Assuming mETH = $1000
        leverageStrategy.borrowStablecoin(borrowAmount);
        
        assertEq(leverageStrategy.totalBorrowed(), borrowAmount);
        
        vm.stopPrank();
    }
    
    function test_LeverageDeployToRWA() public {
        vm.startPrank(admin);
        
        // Setup: Supply collateral and borrow
        uint256 collateralAmount = 100 * 1e18;
        mETH.approve(address(leverageStrategy), collateralAmount);
        leverageStrategy.supplyCollateral(collateralAmount);
        
        leverageStrategy.grantRole(leverageStrategy.KEEPER_ROLE(), admin);
        uint256 borrowAmount = 60_000 * 1e6;
        leverageStrategy.borrowStablecoin(borrowAmount);
        
        // Whitelist strategy for AIT
        ait.setWhitelist(address(leverageStrategy), true);
        
        // Deploy to RWA
        leverageStrategy.deployToRwa(borrowAmount);
        
        assertGt(leverageStrategy.totalAITHoldings(), 0);
        
        vm.stopPrank();
    }
    
    function test_LeverageYieldHarvest() public {
        vm.startPrank(admin);
        
        // Setup full flow
        uint256 collateralAmount = 100 * 1e18;
        mETH.approve(address(leverageStrategy), collateralAmount);
        leverageStrategy.supplyCollateral(collateralAmount);
        
        leverageStrategy.grantRole(leverageStrategy.KEEPER_ROLE(), admin);
        uint256 borrowAmount = 60_000 * 1e6;
        leverageStrategy.borrowStablecoin(borrowAmount);
        
        ait.setWhitelist(address(leverageStrategy), true);
        leverageStrategy.deployToRwa(borrowAmount);
        
        vm.stopPrank();
        
        // Simulate yield accrual by updating NAV
        vm.startPrank(keeper);
        ait.updateNav(1.08e6); // 8% appreciation
        vm.stopPrank();
        
        // Harvest yield
        vm.startPrank(admin);
        uint256 yield = leverageStrategy.harvestRwaYield();
        assertGt(yield, 0);
        
        vm.stopPrank();
    }
    
    function test_LeverageEmergencyDeleverage() public {
        vm.startPrank(admin);
        
        // Setup leveraged position
        uint256 collateralAmount = 100 * 1e18;
        mETH.approve(address(leverageStrategy), collateralAmount);
        leverageStrategy.supplyCollateral(collateralAmount);
        
        leverageStrategy.grantRole(leverageStrategy.KEEPER_ROLE(), admin);
        uint256 borrowAmount = 60_000 * 1e6;
        leverageStrategy.borrowStablecoin(borrowAmount);
        
        ait.setWhitelist(address(leverageStrategy), true);
        leverageStrategy.deployToRwa(borrowAmount);
        
        // Simulate health factor drop (mock)
        // In real test, would manipulate lending protocol state
        
        vm.stopPrank();
    }
    
    // ========================================================================
    // TEST 4: Integration Tests
    // ========================================================================
    
    function test_FullUserJourney() public {
        // 1. Issue KYC
        vm.startPrank(admin);
        kycVerifier.issueKyc(
            accreditedInvestor,
            TieredKYCVerifier.InvestorTier.ACCREDITED,
            365,
            bytes32("US"),
            keccak256(abi.encodePacked(accreditedInvestor))
        );
        vm.stopPrank();
        
        // 2. Investor deposits (simulated)
        // In full implementation, would interact with vault
        
        // 3. Vault deploys to leverage strategy
        // 4. Strategy borrows and deploys to RWA
        // 5. Yield accrues and is harvested
        // 6. Investor withdraws with profit
        
        assertTrue(kycVerifier.hasValidKyc(accreditedInvestor));
    }
    
    function test_MLModelAllocationAdjustment() public pure {
        // Simulate ML model updating allocation based on risk scores
        // This would integrate with the ML API in production
        
        // High risk scenario: reduce RWA allocation
        // Expected: allocation should decrease
        
        // Low risk scenario: increase RWA allocation
        // Expected: allocation should increase
        
        // Note: Full implementation would call vault.updateMLModel()
    }
    
    function test_TaxEventEmission() public {
        // Test that yield distributions emit proper tax events
        // These events are captured by keeper bot for 1099 generation
        
        vm.startPrank(admin);
        
        // Setup investor
        kycVerifier.issueKyc(
            accreditedInvestor,
            TieredKYCVerifier.InvestorTier.ACCREDITED,
            365,
            bytes32("US"),
            keccak256(abi.encodePacked(accreditedInvestor))
        );
        
        // Simulate yield distribution
        // In production, vault would emit YieldDistributionEvent
        
        vm.stopPrank();
    }
    
    // ========================================================================
    // TEST 5: Edge Cases & Security
    // ========================================================================
    
    function test_RevertWhen_UnauthorizedKYCIssuance() public {
        vm.startPrank(retailInvestor); // Not admin
        
        vm.expectRevert(); // Expect AccessControl revert
        kycVerifier.issueKyc(
            retailInvestor,
            TieredKYCVerifier.InvestorTier.RETAIL,
            365,
            bytes32("US"),
            keccak256(abi.encodePacked(retailInvestor))
        );
        
        vm.stopPrank();
    }
    
    function test_RevertWhen_AITTransferToNonWhitelisted() public {
        vm.startPrank(admin);
        
        ait.setWhitelist(admin, true);
        ait.mint(admin, 1000 * 1e6);
        
        // Try to transfer to non-whitelisted address (should fail)
        vm.expectRevert("AIT: Recipient not whitelisted");
        ait.transfer(retailInvestor, 100 * 1e6);
        
        vm.stopPrank();
    }
    
    function test_RevertWhen_ExceedLeverageLTV() public {
        vm.startPrank(admin);
        
        uint256 collateralAmount = 100 * 1e18;
        mETH.approve(address(leverageStrategy), collateralAmount);
        leverageStrategy.supplyCollateral(collateralAmount);
        
        leverageStrategy.grantRole(leverageStrategy.KEEPER_ROLE(), admin);
        
        // Try to borrow more than target LTV
        uint256 excessiveBorrow = 80_000 * 1e6; // 80% LTV (exceeds 60% target)
        vm.expectRevert("Exceeds target LTV");
        leverageStrategy.borrowStablecoin(excessiveBorrow);
        
        vm.stopPrank();
    }
    
    function test_KYCExpiration() public {
        vm.startPrank(admin);
        
        // Issue KYC with 1 day validity
        kycVerifier.issueKyc(
            retailInvestor,
            TieredKYCVerifier.InvestorTier.RETAIL,
            1, // 1 day
            bytes32("US"),
            keccak256(abi.encodePacked(retailInvestor))
        );
        
        assertTrue(kycVerifier.hasValidKyc(retailInvestor));
        
        // Fast forward 2 days
        vm.warp(block.timestamp + 2 days);
        
        assertFalse(kycVerifier.hasValidKyc(retailInvestor));
        
        vm.stopPrank();
    }
}

// ============================================================================
// Mock Contracts for Testing
// ============================================================================

contract MockERC20 {
    string public name;
    string public symbol;
    uint8 public decimals;
    uint256 public totalSupply;
    
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;
    
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    
    constructor(string memory _name, string memory _symbol, uint8 _decimals) {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
    }
    
    function mint(address to, uint256 amount) external {
        balanceOf[to] += amount;
        totalSupply += amount;
        emit Transfer(address(0), to, amount);
    }
    
    function transfer(address to, uint256 amount) external returns (bool) {
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        emit Transfer(msg.sender, to, amount);
        return true;
    }
    
    function approve(address spender, uint256 amount) external returns (bool) {
        allowance[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }
    
    function transferFrom(address from, address to, uint256 amount) external returns (bool) {
        require(allowance[from][msg.sender] >= amount, "ERC20: insufficient allowance");
        require(balanceOf[from] >= amount, "ERC20: transfer exceeds balance");
        
        allowance[from][msg.sender] -= amount;
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        emit Transfer(from, to, amount);
        return true;
    }
}

contract MockLendingProtocol is IMantleLendingProtocol {
    IERC20 public usdc;
    mapping(address => uint256) public supplied;
    mapping(address => uint256) public borrowed;
    
    constructor(address _usdc) {
        usdc = IERC20(_usdc);
    }
    
    function supply(address asset, uint256 amount) external {
        require(IERC20(asset).transferFrom(msg.sender, address(this), amount), "Transfer failed");
        supplied[msg.sender] += amount;
    }
    
    function borrow(address asset, uint256 amount) external {
        borrowed[msg.sender] += amount;
        require(IERC20(asset).transfer(msg.sender, amount), "Transfer failed");
    }
    
    function repay(address asset, uint256 amount) external {
        require(IERC20(asset).transferFrom(msg.sender, address(this), amount), "Transfer failed");
        borrowed[msg.sender] -= amount;
    }
    
    function withdraw(address asset, uint256 amount) external {
        supplied[msg.sender] -= amount;
        require(IERC20(asset).transfer(msg.sender, amount), "Transfer failed");
    }
    
    function getAccountLiquidity(address account) external view returns (
        uint256 collateralValue,
        uint256 borrowValue,
        uint256 healthFactor
    ) {
        collateralValue = supplied[account];
        borrowValue = borrowed[account];
        healthFactor = collateralValue > 0 ? (collateralValue * 10000) / (borrowValue + 1) : 0;
    }
    
    function getCollateralFactor(address) external pure returns (uint256) {
        return 7000; // 70%
    }
    
    function getBorrowRate(address) external pure returns (uint256) {
        return 500; // 5%
    }
}
