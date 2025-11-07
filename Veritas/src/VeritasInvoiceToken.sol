// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

/**
 * @title VeritasInvoiceToken (AIT)
 * @notice Synthetic asset representing tokenized claims on future cash flows
 * @dev ERC-20 token backed by off-chain corporate invoice pool
 * 
 * KEY DISTINCTION:
 * - NOT ownership of the invoices themselves
 * - Represents CLAIM on future cash flows when invoices are paid
 * - Synthetic asset with oracle-based pricing
 * - Compliant with securities regulations (Reg D/Reg S)
 */
contract VeritasInvoiceToken is IERC20, AccessControl {
    
    bytes32 public constant ISSUER_ROLE = keccak256("ISSUER_ROLE");
    bytes32 public constant ORACLE_ROLE = keccak256("ORACLE_ROLE");
    
    string public constant name = "Veritas Invoice Token";
    string public constant symbol = "AIT";
    uint8 public constant decimals = 6; // USDC-compatible
    
    uint256 private _totalSupply;
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;
    
    // Synthetic Asset Metadata
    struct UnderlyingPool {
        bytes32 poolId;              // e.g., "CORP-INV-Q1-2025"
        uint256 totalFaceValue;      // Sum of all invoice values
        uint256 numberOfInvoices;    // Number of invoices in pool
        uint256 weightedMaturity;    // Weighted avg maturity (days)
        uint256 expectedYield;       // Expected APY in basis points
        uint256 realizedYield;       // Actual yield collected
        uint256 defaultRate;         // Current default rate (bps)
    }
    
    UnderlyingPool public pool;
    
    // Oracle-based pricing (NAV per token)
    uint256 public navPerToken; // Net Asset Value in USDC (scaled by 1e6)
    uint256 public lastNavUpdate;
    
    // Compliance: Transfer restrictions
    mapping(address => bool) public whitelisted;
    
    event PoolInitialized(bytes32 indexed poolId, uint256 totalFaceValue);
    event NAVUpdated(uint256 oldNav, uint256 newNav, uint256 timestamp);
    event CashFlowDistributed(uint256 amount, uint256 invoicesPaid);
    event DefaultRecorded(uint256 defaultAmount, uint256 invoiceId);
    
    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ISSUER_ROLE, msg.sender);
        _grantRole(ORACLE_ROLE, msg.sender);
        
        navPerToken = 1e6; // Start at $1.00
        lastNavUpdate = block.timestamp;
    }
    
    /**
     * @notice Initialize underlying invoice pool
     * @param poolId Unique identifier for invoice pool
     * @param totalFaceValue Total value of invoices (USDC)
     * @param numberOfInvoices Count of invoices
     * @param weightedMaturity Average maturity in days
     * @param expectedYield Expected APY in basis points
     */
    function initializePool(
        bytes32 poolId,
        uint256 totalFaceValue,
        uint256 numberOfInvoices,
        uint256 weightedMaturity,
        uint256 expectedYield
    ) external onlyRole(ISSUER_ROLE) {
        pool = UnderlyingPool({
            poolId: poolId,
            totalFaceValue: totalFaceValue,
            numberOfInvoices: numberOfInvoices,
            weightedMaturity: weightedMaturity,
            expectedYield: expectedYield,
            realizedYield: 0,
            defaultRate: 0
        });
        
        emit PoolInitialized(poolId, totalFaceValue);
    }
    
    /**
     * @notice Update NAV based on oracle price feed
     * @dev Called by oracle after invoice payments/defaults
     * @param newNav New net asset value per token
     */
    function updateNav(uint256 newNav) external onlyRole(ORACLE_ROLE) {
        uint256 oldNav = navPerToken;
        navPerToken = newNav;
        lastNavUpdate = block.timestamp;
        
        emit NAVUpdated(oldNav, newNav, block.timestamp);
    }
    
    /**
     * @notice Record cash flow distribution (invoices paid)
     * @param amount USDC amount collected
     * @param invoicesPaid Number of invoices that paid
     */
    function recordCashFlow(uint256 amount, uint256 invoicesPaid) 
        external 
        onlyRole(ORACLE_ROLE) 
    {
        pool.realizedYield += amount;
        
        // Update NAV based on collected cash
        if (_totalSupply > 0) {
            navPerToken = (pool.totalFaceValue + pool.realizedYield) * 1e6 / _totalSupply;
        }
        
        emit CashFlowDistributed(amount, invoicesPaid);
    }
    
    /**
     * @notice Record invoice default
     * @param defaultAmount USDC amount defaulted
     * @param invoiceId Identifier of defaulted invoice
     */
    function recordDefault(uint256 defaultAmount, uint256 invoiceId) 
        external 
        onlyRole(ORACLE_ROLE) 
    {
        pool.totalFaceValue -= defaultAmount;
        pool.defaultRate = (defaultAmount * 10000) / pool.totalFaceValue;
        
        // Reduce NAV due to default
        if (_totalSupply > 0) {
            navPerToken = pool.totalFaceValue * 1e6 / _totalSupply;
        }
        
        emit DefaultRecorded(defaultAmount, invoiceId);
    }
    
    /**
     * @notice Mint AIT tokens (issuer only)
     * @dev Represents new claims on invoice cash flows
     */
    function mint(address to, uint256 amount) external onlyRole(ISSUER_ROLE) {
        require(whitelisted[to], "AIT: Recipient not whitelisted");
        _mint(to, amount);
    }
    
    /**
     * @notice Burn AIT tokens on redemption
     */
    function burn(address from, uint256 amount) external onlyRole(ISSUER_ROLE) {
        _burn(from, amount);
    }
    
    /**
     * @notice Whitelist address for compliance
     */
    function setWhitelist(address account, bool status) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        whitelisted[account] = status;
    }
    
    // ============ ERC-20 Implementation with Restrictions ============
    
    function totalSupply() external view override returns (uint256) {
        return _totalSupply;
    }
    
    function balanceOf(address account) external view override returns (uint256) {
        return _balances[account];
    }
    
    function transfer(address to, uint256 amount) external override returns (bool) {
        require(whitelisted[to], "AIT: Recipient not whitelisted");
        _transfer(msg.sender, to, amount);
        return true;
    }
    
    function allowance(address owner, address spender) 
        external 
        view 
        override 
        returns (uint256) 
    {
        return _allowances[owner][spender];
    }
    
    function approve(address spender, uint256 amount) 
        external 
        override 
        returns (bool) 
    {
        _approve(msg.sender, spender, amount);
        return true;
    }
    
    function transferFrom(address from, address to, uint256 amount) 
        external 
        override 
        returns (bool) 
    {
        require(whitelisted[to], "AIT: Recipient not whitelisted");
        
        uint256 currentAllowance = _allowances[from][msg.sender];
        require(currentAllowance >= amount, "ERC20: insufficient allowance");
        
        _approve(from, msg.sender, currentAllowance - amount);
        _transfer(from, to, amount);
        
        return true;
    }
    
    function _mint(address account, uint256 amount) internal {
        _totalSupply += amount;
        _balances[account] += amount;
        emit Transfer(address(0), account, amount);
    }
    
    function _burn(address account, uint256 amount) internal {
        require(_balances[account] >= amount, "ERC20: burn exceeds balance");
        _balances[account] -= amount;
        _totalSupply -= amount;
        emit Transfer(account, address(0), amount);
    }
    
    function _transfer(address from, address to, uint256 amount) internal {
        require(_balances[from] >= amount, "ERC20: transfer exceeds balance");
        _balances[from] -= amount;
        _balances[to] += amount;
        emit Transfer(from, to, amount);
    }
    
    function _approve(address owner, address spender, uint256 amount) internal {
        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }
}
