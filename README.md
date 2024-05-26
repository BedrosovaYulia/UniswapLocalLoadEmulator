# Uniswap Load Emulator

This project is a load emulator for various versions of Uniswap exchanges on a local Ethereum fork. It generates random transactions for testing other local applications.

## Requirements

- Docker
- Docker Compose

## Installation

1. Clone the repository:

```sh
git clone https://github.com/yourusername/uniswap-load-emulator.git
cd uniswap-load-emulator
```

2. Create a `.env` file in the root directory of the project and add the following environment variables:

```env
RPC_URL=http://host.docker.internal:8545/
PRIVATE_KEY=0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
```

Replace `PRIVATE_KEY` with your private key.

## Usage

1. Ensure that your local Ethereum fork is running and accessible at the address specified in `RPC_URL`.

2. Start Docker Compose:

```sh
docker-compose up --build
```

This will create and start the `uniswap-load-emulator` container, which will generate random transactions for Uniswap.