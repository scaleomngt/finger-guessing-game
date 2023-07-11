#  finger-guessing-game

 finger-guessing-game

The player's data while playing the game is sensitive, We need to utilize zero-knowledge proofs to ensure that sensitive data remains encrypted and under the control of the data owner at all times. Aleo ZKML addresses this issue by integrating privacy into the smart contract layer. For third parties that require verification of user data, ZKP (Zero-Knowledge Proof) smart contracts can be written to prove the correctness of computations without revealing the underlying data.

Program Description

Program URLs:
* 1. https://github.com/scaleomngt/finger-guessing-game/tree/main/front/Stone-cloth-scissors --Frontend code
* 2. https://github.com/scaleomngt/finger-guessing-game --Server code
* 3. https://github.com/scaleomngt/finger-guessing-game/tree/main/aleo/game_2ase7z  --Leo code

Installation Instructions

Frontend Deployment

(1) Step 1: Download the dependencies required for the program
   
    npm install

(2) Step 2: Verify if all dependencies are successfully downloaded by running the program locally

    npm run serve 

(3) Step 3: Access the configured IP address and port in a browser to use the application

Backend Deployment
    go build

leoDeployment
- `cd aleo/game_2ase7z && leo build`
  snarkos developer deploy "game_2ase7z.aleo" --private-key "APrivateKey11111111"  --query "https://vm.aleo.org/api"  --path "game_2ase7z/build/"  --broadcast "https://vm.aleo.org/api/testnet3/transaction/broadcast"  --fee 400000  --record "$RECORD"


