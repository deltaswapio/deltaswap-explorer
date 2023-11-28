import { ChainName, coalesceChainId } from '@deltaswapio/deltaswap-sdk/lib/cjs/utils/consts';
import { readFileSync, writeFileSync } from 'fs';
import { env } from '../config';
import BaseDB from './BaseDB';
import { WHTransaction, WHTransferRedeemed } from './types';

const ENCODING = 'utf8';
const DELTASWAP_TX_FILE: string = env.JSON_WH_TXS_FILE;
const GLOBAL_TX_FILE: string = env.JSON_GLOBAL_TXS_FILE;
const DELTASWAP_LAST_BLOCKS_FILE: string = env.JSON_LAST_BLOCKS_FILE;

export default class JsonDB extends BaseDB {
  deltaswapTxFile: WHTransaction[] = [];
  redeemedTxFile: WHTransferRedeemed[] = [];

  constructor() {
    super('JsonDB');
    this.deltaswapTxFile = [];
    this.redeemedTxFile = [];
    this.lastBlocksByChain = [];
    this.logger.info('Connecting...');
  }

  async connect(): Promise<void> {
    try {
      const whTxsFileRawData = readFileSync(DELTASWAP_TX_FILE, ENCODING);
      this.deltaswapTxFile = whTxsFileRawData ? JSON.parse(whTxsFileRawData) : [];
      this.logger.info(`${DELTASWAP_TX_FILE} file ready`);
    } catch (e) {
      this.logger.warn(`${DELTASWAP_TX_FILE} file does not exists, creating new file`);
      this.deltaswapTxFile = [];
    }

    try {
      const whRedeemedTxsFileRawData = readFileSync(GLOBAL_TX_FILE, ENCODING);
      this.redeemedTxFile = whRedeemedTxsFileRawData ? JSON.parse(whRedeemedTxsFileRawData) : [];
      this.logger.info(`${GLOBAL_TX_FILE} file ready`);
    } catch (e) {
      this.logger.warn(`${GLOBAL_TX_FILE} file does not exists, creating new file`);
      this.redeemedTxFile = [];
    }
  }

  async disconnect(): Promise<void> {
    this.logger.info('Disconnecting...');
    this.logger.info('Disconnected');
  }

  async isConnected() {
    return true;
  }

  async getLastBlocksProcessed(): Promise<void> {
    try {
      const lastBlocksByChain = readFileSync(DELTASWAP_LAST_BLOCKS_FILE, ENCODING);
      this.lastBlocksByChain = lastBlocksByChain ? JSON.parse(lastBlocksByChain) : [];
      this.logger.info(`${DELTASWAP_LAST_BLOCKS_FILE} file ready`);
    } catch (e) {
      this.logger.warn(`${DELTASWAP_LAST_BLOCKS_FILE} file does not exists, creating new file`);
      this.lastBlocksByChain = [];
    }
  }

  async storeWhTxs(chainName: ChainName, whTxs: WHTransaction[]): Promise<void> {
    try {
      for (let i = 0; i < whTxs.length; i++) {
        let message = 'Insert Deltaswap Transaction Event Log to JSON file';
        const currentWhTx = whTxs[i];
        const { id } = currentWhTx;

        currentWhTx.eventLog.unsignedVaa = Buffer.isBuffer(currentWhTx.eventLog.unsignedVaa)
          ? Buffer.from(currentWhTx.eventLog.unsignedVaa).toString('base64')
          : currentWhTx.eventLog.unsignedVaa;

        const whTxIndex = this.deltaswapTxFile?.findIndex((whTx) => whTx.id === id.toString());

        if (whTxIndex >= 0) {
          const whTx = this.deltaswapTxFile[whTxIndex];

          whTx.eventLog.updatedAt = new Date();
          whTx.eventLog.revision ? (whTx.eventLog.revision += 1) : (whTx.eventLog.revision = 1);

          message = 'Update Deltaswap Transaction Event Log to JSON file';
        } else {
          this.deltaswapTxFile.push(currentWhTx);
        }

        writeFileSync(DELTASWAP_TX_FILE, JSON.stringify(this.deltaswapTxFile, null, 2), ENCODING);

        if (currentWhTx) {
          const { id, eventLog } = currentWhTx;
          const { blockNumber, txHash, emitterChain } = eventLog;

          this.logger.info({
            id,
            blockNumber,
            chainName,
            txHash,
            emitterChain,
            message,
          });
        }
      }
    } catch (e: unknown) {
      this.logger.error(`Error Upsert Deltaswap Transaction Event Log: ${e}`);
    }
  }

  async storeRedeemedTxs(chainName: ChainName, redeemedTxs: WHTransferRedeemed[]): Promise<void> {
    // For JsonDB we are only pushing all the "redeemed" logs into GLOBAL_TX_FILE simulating a globalTransactions collection

    try {
      for (let i = 0; i < redeemedTxs.length; i++) {
        let message = 'Insert Deltaswap Transfer Redeemed Event Log to JSON file';
        const currentRedeemedTx = redeemedTxs[i];
        const { id, destinationTx } = currentRedeemedTx;
        const { method, status } = destinationTx;

        const whTxIndex = this.deltaswapTxFile?.findIndex((whTx) => whTx.id === id.toString());

        if (whTxIndex >= 0) {
          const whTx = this.deltaswapTxFile[whTxIndex];

          whTx.status = status;
          whTx.eventLog.updatedAt = new Date();
          whTx.eventLog.revision ? (whTx.eventLog.revision += 1) : (whTx.eventLog.revision = 1);

          writeFileSync(DELTASWAP_TX_FILE, JSON.stringify(this.deltaswapTxFile, null, 2), ENCODING);
        }

        const whRedeemedTxIndex = this.redeemedTxFile?.findIndex(
          (whRedeemedTx) => whRedeemedTx.id === id.toString(),
        );

        if (whRedeemedTxIndex >= 0) {
          const whRedeemedTx = this.redeemedTxFile[whRedeemedTxIndex];

          whRedeemedTx.destinationTx.method = method;
          whRedeemedTx.destinationTx.status = status;
          whRedeemedTx.destinationTx.updatedAt = new Date();
          whRedeemedTx.revision ? (whRedeemedTx.revision += 1) : (whRedeemedTx.revision = 1);

          message = 'Update Deltaswap Transfer Redeemed Event Log to JSON file';
        } else {
          this.redeemedTxFile.push(currentRedeemedTx);
        }

        writeFileSync(GLOBAL_TX_FILE, JSON.stringify(this.redeemedTxFile, null, 2), ENCODING);

        if (currentRedeemedTx) {
          const { id, destinationTx } = currentRedeemedTx;
          const { chainId } = destinationTx;

          this.logger.info({
            id,
            chainId,
            chainName,
            message,
          });
        }
      }
    } catch (e: unknown) {
      this.logger.error(`Error Upsert Deltaswap Transfer Redeemed Event Log: ${e}`);
    }
  }

  async storeLatestProcessBlock(
    chain: ChainName,
    lastBlock: number,
    lastSequenceNumber: number | null,
  ): Promise<void> {
    const chainId = coalesceChainId(chain);
    const updatedLastBlocksByChain = [...this.lastBlocksByChain];
    const itemIndex = updatedLastBlocksByChain.findIndex((item) => {
      if ('id' in item) return item.id === chain;
      return false;
    });

    if (itemIndex >= 0) {
      updatedLastBlocksByChain[itemIndex] = {
        ...updatedLastBlocksByChain[itemIndex],
        blockNumber: lastBlock,
        lastSequenceNumber,
        updatedAt: new Date(),
      };
    } else {
      updatedLastBlocksByChain.push({
        id: chain,
        blockNumber: lastBlock,
        lastSequenceNumber,
        chainId,
        createdAt: new Date(),
        updatedAt: new Date(),
      });
    }

    this.lastBlocksByChain = updatedLastBlocksByChain;

    try {
      writeFileSync(
        DELTASWAP_LAST_BLOCKS_FILE,
        JSON.stringify(this.lastBlocksByChain, null, 2),
        ENCODING,
      );
    } catch (e: unknown) {
      this.logger.error(`Error Insert latest processed block: ${e}`);
    }
  }
}
