// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
/**
 * @title TieredKYCVerifier
 * @notice Enhanced SBT verifier with investor tiers and caps
 * @dev Implements different investment limits based on accreditation level
 */
contract TieredKYCVerifier is AccessControl {
    
    bytes32 public constant KYC_ISSUER = keccak256("KYC_ISSUER");
    
    enum InvestorTier {
        NONE,           // 0 - Not verified
        RETAIL,         // 1 - Retail investor ($10k cap)
        ACCREDITED,     // 2 - Accredited investor ($500k cap)
        QUALIFIED,      // 3 - Qualified purchaser ($5M cap)
        INSTITUTIONAL   // 4 - Institutional (unlimited)
    }
    
    struct KYCProfile {
        InvestorTier tier;
        uint256 investmentCap;  // Max investment in USDC
        uint256 currentInvestment; // Current position
        uint256 issuedAt;
        uint256 expiresAt;
        bool revoked;
        bytes32 jurisdiction;
        bytes32 kycIdHash; // Privacy-preserving identifier
    }
    
    mapping(address => KYCProfile) public kycProfiles;
    
    // Default caps per tier (in USDC, scaled by 1e6)
    mapping(InvestorTier => uint256) public tierCaps;
    
    event KYCIssued(address indexed investor, InvestorTier tier, uint256 cap);
    event InvestmentRecorded(address indexed investor, uint256 amount, uint256 newTotal);
    event TierUpgraded(address indexed investor, InvestorTier oldTier, InvestorTier newTier);
    event KYCRevoked(address indexed investor, string reason);
    
    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(KYC_ISSUER, msg.sender);
        
        // Set default caps
        tierCaps[InvestorTier.RETAIL] = 10_000 * 1e6;      // $10k
        tierCaps[InvestorTier.ACCREDITED] = 500_000 * 1e6; // $500k
        tierCaps[InvestorTier.QUALIFIED] = 5_000_000 * 1e6; // $5M
        tierCaps[InvestorTier.INSTITUTIONAL] = type(uint256).max; // Unlimited
    }
    
    /**
     * @notice Issue tiered KYC to investor
     */
    function issueKYC(
        address investor,
        InvestorTier tier,
        uint256 validityDays,
        bytes32 jurisdiction,
        bytes32 kycIdHash
    ) external onlyRole(KYC_ISSUER) {
        require(tier != InvestorTier.NONE, "Invalid tier");
        
        kycProfiles[investor] = KYCProfile({
            tier: tier,
            investmentCap: tierCaps[tier],
            currentInvestment: 0,
            issuedAt: block.timestamp,
            expiresAt: block.timestamp + (validityDays * 1 days),
            revoked: false,
            jurisdiction: jurisdiction,
            kycIdHash: kycIdHash
        });
        
        emit KYCIssued(investor, tier, tierCaps[tier]);
    }
    
    /**
     * @notice Check if investor can invest amount
     */
    function canInvest(address investor, uint256 amount) 
        external 
        view 
        returns (bool, string memory reason) 
    {
        KYCProfile memory profile = kycProfiles[investor];
        
        if (profile.tier == InvestorTier.NONE) {
            return (false, "No KYC");
        }
        if (profile.revoked) {
            return (false, "KYC revoked");
        }
        if (block.timestamp > profile.expiresAt) {
            return (false, "KYC expired");
        }
        if (profile.currentInvestment + amount > profile.investmentCap) {
            return (false, "Exceeds investment cap");
        }
        
        return (true, "");
    }
    
    /**
     * @notice Record investment (called by vault)
     */
    function recordInvestment(address investor, uint256 amount) external {
        KYCProfile storage profile = kycProfiles[investor];
        require(profile.currentInvestment + amount <= profile.investmentCap, "Exceeds cap");
        
        profile.currentInvestment += amount;
        
        emit InvestmentRecorded(investor, amount, profile.currentInvestment);
    }
    
    /**
     * @notice Upgrade investor tier
     */
    function upgradeTier(address investor, InvestorTier newTier) 
        external 
        onlyRole(KYC_ISSUER) 
    {
        KYCProfile storage profile = kycProfiles[investor];
        require(newTier > profile.tier, "Not an upgrade");
        
        InvestorTier oldTier = profile.tier;
        profile.tier = newTier;
        profile.investmentCap = tierCaps[newTier];
        
        emit TierUpgraded(investor, oldTier, newTier);
    }
    
    /**
     * @notice Revoke KYC
     */
    function revokeKYC(address investor, string calldata reason) 
        external 
        onlyRole(KYC_ISSUER) 
    {
        kycProfiles[investor].revoked = true;
        emit KYCRevoked(investor, reason);
    }
    
    /**
     * @notice Check if investor has valid KYC
     */
    function hasValidKYC(address investor) external view returns (bool) {
        KYCProfile memory profile = kycProfiles[investor];
        return profile.tier != InvestorTier.NONE 
            && !profile.revoked 
            && block.timestamp <= profile.expiresAt;
    }
    
    /**
     * @notice Get investor tier
     */
    function getTier(address investor) external view returns (InvestorTier) {
        return kycProfiles[investor].tier;
    }
    
    /**
     * @notice Get remaining investment capacity
     */
    function getRemainingCapacity(address investor) external view returns (uint256) {
        KYCProfile memory profile = kycProfiles[investor];
        if (profile.currentInvestment >= profile.investmentCap) {
            return 0;
        }
        return profile.investmentCap - profile.currentInvestment;
    }
    
    /**
     * @notice Update tier caps (governance)
     */
    function updateTierCap(InvestorTier tier, uint256 newCap) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        tierCaps[tier] = newCap;
    }
}
