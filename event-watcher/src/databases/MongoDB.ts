import { ChainName, coalesceChainId } from '@deltaswapio/deltaswap-sdk/lib/cjs/utils/consts';
import BaseDB from './BaseDB';
import { LastBlockByChain, WHTransaction, WHTransferRedeemed } from './types';
import * as mongoDB from 'mongodb';
import { env } from '../config';

const DELTASWAP_TX_COLLECTION: string = 'deltaswapTxs';
const GLOBAL_TX_COLLECTION: string = 'globalTransactions';
const DELTASWAP_LAST_BLOCK_COLLECTION: string = 'lastBlocksByChain';

export default class MongoDB extends BaseDB {
  private client: mongoDB.MongoClient | null = null;
  private db: mongoDB.Db | null = null;
  private deltaswapTxCollection: mongoDB.Collection | null = null;
  private globalTxCollection: mongoDB.Collection | null = null;
  private lastTxBlockByChainCollection: mongoDB.Collection | null = null;

  constructor() {
    super('MongoDB');
    this.logger.info('Connecting...');
    this.client = new mongoDB.MongoClient(env.MONGODB_URI as string);
    this.db = this.client.db(env.MONGODB_DATABASE ?? 'deltaswap');
    this.deltaswapTxCollection = this.db.collection(DELTASWAP_TX_COLLECTION);
    this.globalTxCollection = this.db.collection(GLOBAL_TX_COLLECTION);
    this.lastTxBlockByChainCollection = this.db.collection(DELTASWAP_LAST_BLOCK_COLLECTION);
  }

  async connect(): Promise<void> {
    try {
      await this.client?.connect();

      this.logger.info('Connected');
    } catch (e) {
      throw new Error(`[MongoDB] Error: ${e}`);
    }
  }

  async disconnect(): Promise<void> {
    this.logger.info('Disconnecting...');
    await this.client?.close();
    this.logger.info('Disconnected');
  }

  async isConnected() {
    try {
      await this.db?.command({ ping: 1 });
      return true;
    } catch (error: unknown) {
      return false;
    }
  }

  async getLastBlocksProcessed(): Promise<void> {
    try {
      const lastBlocksByChain = await this.lastTxBlockByChainCollection?.find().toArray();
      this.lastBlocksByChain = (lastBlocksByChain as unknown as LastBlockByChain[]) || [];
    } catch (error: unknown) {
      this.logger.warn(`Error while getting last blocks processed: ${error}`);
      this.lastBlocksByChain = [];
    }
  }

  async storeWhTxs(chainName: ChainName, whTxs: WHTransaction[]): Promise<void> {
    try {
      for (let i = 0; i < whTxs.length; i++) {
        let message = `Insert Deltaswap Transaction Event Log to ${DELTASWAP_TX_COLLECTION} collection`;
        const currentWhTx = whTxs[i];
        const { id, ...rest } = currentWhTx;

        const whTxDocument = await this.deltaswapTxCollection?.findOne({
          _id: id as unknown as mongoDB.ObjectId,
        });

        if (whTxDocument) {
          await this.deltaswapTxCollection?.findOneAndUpdate(
            { _id: id as unknown as mongoDB.ObjectId },
            {
              $set: { 'eventLog.updatedAt': new Date() },
              $inc: { 'eventLog.revision': 1 },
            },
            { returnDocument: 'after' },
          );

          message = `Update Deltaswap Transaction Event Log to ${DELTASWAP_TX_COLLECTION} collection`;
        } else {
          await this.deltaswapTxCollection?.insertOne({
            _id: id as unknown as mongoDB.ObjectId,
            ...rest,
          });
        }

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
    try {
      for (let i = 0; i < redeemedTxs.length; i++) {
        let message = `Insert Deltaswap Transfer Redeemed Event Log to ${GLOBAL_TX_COLLECTION} collection`;
        const currentWhRedeemedTx = redeemedTxs[i];
        const { id, destinationTx, ...rest } = currentWhRedeemedTx;
        const { status } = destinationTx;

        const whTxResponse = await this.deltaswapTxCollection?.findOneAndUpdate(
          { _id: id as unknown as mongoDB.ObjectId },
          {
            $set: {
              'eventLog.updatedAt': new Date(),
              status: status,
            },
            $inc: { 'eventLog.revision': 1 },
          },
          { returnDocument: 'after' },
        );

        if (!whTxResponse?.value) {
          this.logger.info(
            `Error Update Deltaswap Transfer Redeemed Event Log: ${id} does not exist on ${DELTASWAP_TX_COLLECTION} collection`,
          );
        }

        const globalTxDocument = await this.globalTxCollection?.findOne({
          _id: id as unknown as mongoDB.ObjectId,
        });

        if (globalTxDocument) {
          message = `Update Deltaswap Transfer Redeemed Event Log to ${GLOBAL_TX_COLLECTION} collection`;
          const { destinationTx: globalTxDocumentDestinationTx } = globalTxDocument;

          if (!globalTxDocumentDestinationTx) {
            await this.globalTxCollection?.findOneAndUpdate(
              { _id: id as unknown as mongoDB.ObjectId },
              {
                $set: { destinationTx },
                $inc: { revision: 1 },
              },
              { returnDocument: 'after' },
            );
          } else {
            message = `Already exists Deltaswap Transfer Redeemed Event Log on ${GLOBAL_TX_COLLECTION} collection`;
          }
        } else {
          await this.globalTxCollection?.insertOne({
            _id: id as unknown as mongoDB.ObjectId,
            destinationTx,
            ...rest,
          });
        }

        if (currentWhRedeemedTx) {
          const { id, destinationTx } = currentWhRedeemedTx;
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
      this.logger.error(`Error Update Deltaswap Transfer Redeemed Event Log: ${e}`);
    }
  }

  async storeLatestProcessBlock(
    chain: ChainName,
    lastBlock: number,
    lastSequenceNumber: number | null,
  ): Promise<void> {
    const chainId = coalesceChainId(chain);

    try {
      await this.lastTxBlockByChainCollection?.findOneAndUpdate(
        { _id: chain as unknown as mongoDB.ObjectId },
        {
          $setOnInsert: {
            chainId,
            createdAt: new Date(),
          },
          $set: {
            blockNumber: lastBlock,
            lastSequenceNumber,
            updatedAt: new Date(),
          },
        },
        { upsert: true },
      );
    } catch (e: unknown) {
      this.logger.error(`Error while storing latest processed block: ${e}`);
    }
  }
}
