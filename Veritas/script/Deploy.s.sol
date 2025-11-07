// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/VeritasInvoiceToken.sol";
import "../src/LeveragedRWAStrategy.sol";
import "../src/TieredKYCVerifier.sol";

/**
 * @title DeployVeritas
 * @notice Deployment script for Veritas RWA Vault system
 * @dev Run with: forge script script/Deploy.s.sol:DeployVeritas --rpc-url $MANTLE_RPC --broadcast
 */
contract DeployVeritas is Script {
    
    // Mantle Sepolia addresses (update for mainnet)
    address constant METH = 0x0000000000000000000000000000000000000000; // Update with actual mETH
    address constant USDC = 0x0000000000000000000000000000000000000000; // Update with actual USDC
    address constant LENDING_PROTOCOL = 0x0000000000000000000000000000000000000000; // INIT Capital or Lendle
    
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        
        console.log("Deploying Veritas RWA Vault System...");
        console.log("Deployer:", deployer);
        console.log("Chain ID:", block.chainid);
        
        vm.startBroadcast(deployerPrivateKey);
        
        // 1. Deploy VeritasInvoiceToken (AIT)
        console.log("\n1. Deploying VeritasInvoiceToken...");
        VeritasInvoiceToken ait = new VeritasInvoiceToken();
        console.log("AIT deployed at:", address(ait));
        
        // 2. Deploy TieredKYCVerifier
        console.log("\n2. Deploying TieredKYCVerifier...");
        TieredKYCVerifier kycVerifier = new TieredKYCVerifier();
        console.log("KYC Verifier deployed at:", address(kycVerifier));
        
        // 3. Deploy LeveragedRWAStrategy
        console.log("\n3. Deploying LeveragedRWAStrategy...");
        LeveragedRWAStrategy leverageStrategy = new LeveragedRWAStrategy(
            METH,
            USDC,
            LENDING_PROTOCOL,
            address(ait)
        );
        console.log("Leverage Strategy deployed at:", address(leverageStrategy));
        
        // 4. Setup initial configuration
        console.log("\n4. Configuring contracts...");
        
        // Whitelist leverage strategy for AIT minting
        ait.setWhitelist(address(leverageStrategy), true);
        console.log("- Leverage strategy whitelisted for AIT");
        
        // Grant oracle role to deployer (will be transferred to keeper bot)
        ait.grantRole(ait.ORACLE_ROLE(), deployer);
        console.log("- Oracle role granted to deployer");
        
        // Grant keeper role to deployer (will be transferred to keeper bot)
        leverageStrategy.grantRole(leverageStrategy.KEEPER_ROLE(), deployer);
        console.log("- Keeper role granted to deployer");
        
        // Initialize sample invoice pool
        ait.initializePool(
            keccak256("Veritas-DEMO-POOL-2025"),
            1_000_000 * 1e6, // $1M face value
            50,              // 50 invoices
            90,              // 90 day weighted maturity
            800              // 8% expected APY
        );
        console.log("- Sample invoice pool initialized");
        
        vm.stopBroadcast();
        
        // 5. Print deployment summary
        console.log("\n========================================");
        console.log("DEPLOYMENT SUMMARY");
        console.log("========================================");
        console.log("VeritasInvoiceToken:", address(ait));
        console.log("TieredKYCVerifier:", address(kycVerifier));
        console.log("LeveragedRWAStrategy:", address(leverageStrategy));
        console.log("========================================");
        console.log("\nNext steps:");
        console.log("1. Update keeper bot config with contract addresses");
        console.log("2. Transfer oracle role to keeper bot");
        console.log("3. Transfer keeper role to keeper bot");
        console.log("4. Verify contracts on Mantle Explorer");
        console.log("========================================\n");
        
        // 6. Save deployment addresses to file
        string memory deploymentInfo = string(abi.encodePacked(
            "{\n",
            '  "network": "mantle-sepolia",\n',
            '  "chainId": ', vm.toString(block.chainid), ',\n',
            '  "deployer": "', vm.toString(deployer), '",\n',
            '  "contracts": {\n',
            '    "VeritasInvoiceToken": "', vm.toString(address(ait)), '",\n',
            '    "TieredKYCVerifier": "', vm.toString(address(kycVerifier)), '",\n',
            '    "LeveragedRWAStrategy": "', vm.toString(address(leverageStrategy)), '"\n',
            '  }\n',
            '}\n'
        ));
        
        vm.writeFile("deployment.json", deploymentInfo);
        console.log("Deployment info saved to deployment.json");
    }
}

/**
 * @title SetupDemo
 * @notice Setup script for demo/testing
 */
contract SetupDemo is Script {
    
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        // Load deployment addresses
        address aitAddress = vm.envAddress("AIT_ADDRESS");
        address kycAddress = vm.envAddress("KYC_ADDRESS");
        
        vm.startBroadcast(deployerPrivateKey);
        
        VeritasInvoiceToken ait = VeritasInvoiceToken(aitAddress);
        TieredKYCVerifier kyc = TieredKYCVerifier(kycAddress);
        
        // Create demo investors
        address demoRetail = 0x1111111111111111111111111111111111111111;
        address demoAccredited = 0x2222222222222222222222222222222222222222;
        address demoInstitutional = 0x3333333333333333333333333333333333333333;
        
        console.log("Setting up demo investors...");
        
        // Issue KYC for demo investors
        kyc.issueKyc(
            demoRetail,
            TieredKYCVerifier.InvestorTier.RETAIL,
            365,
            bytes32("US"),
            keccak256(abi.encodePacked(demoRetail))
        );
        console.log("- Retail investor KYC issued:", demoRetail);
        
        kyc.issueKyc(
            demoAccredited,
            TieredKYCVerifier.InvestorTier.ACCREDITED,
            365,
            bytes32("US"),
            keccak256(abi.encodePacked(demoAccredited))
        );
        console.log("- Accredited investor KYC issued:", demoAccredited);
        
        kyc.issueKyc(
            demoInstitutional,
            TieredKYCVerifier.InvestorTier.INSTITUTIONAL,
            365,
            bytes32("US"),
            keccak256(abi.encodePacked(demoInstitutional))
        );
        console.log("- Institutional investor KYC issued:", demoInstitutional);
        
        // Whitelist demo investors for AIT
        ait.setWhitelist(demoRetail, true);
        ait.setWhitelist(demoAccredited, true);
        ait.setWhitelist(demoInstitutional, true);
        console.log("- Demo investors whitelisted for AIT");
        
        vm.stopBroadcast();
        
        console.log("\nDemo setup complete!");
    }
}
