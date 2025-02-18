import { CHAIN_ID_SOLANA } from '@deltaswapio/deltaswap-sdk/lib/cjs/utils/consts';
import { expect, test } from '@jest/globals';
import { INITIAL_DEPLOYMENT_BLOCK_BY_CHAIN } from '../../common';
import { getDB } from '../utils';

test('getResumeBlockByChain', async () => {
  const db = getDB();
  const fauxBlock = 98765;
  db.lastBlocksByChain = [
    {
      id: 'solana',
      blockNumber: fauxBlock,
      chainId: CHAIN_ID_SOLANA,
      lastSequenceNumber: 0,
      createdAt: new Date(),
      updatedAt: new Date(),
    },
  ];
  // if a chain is in the database, that number should be returned
  expect(await db.getLastBlockByChain('solana')).toEqual(fauxBlock);
  expect(await db.getResumeBlockByChain('solana')).toEqual(Number(fauxBlock) + 1);
  // if a chain is not in the database, the initial deployment block should be returned
  expect(INITIAL_DEPLOYMENT_BLOCK_BY_CHAIN.moonbeam).toBeDefined();
  expect(await db.getResumeBlockByChain('moonbeam')).toEqual(
    Number(INITIAL_DEPLOYMENT_BLOCK_BY_CHAIN.moonbeam),
  );
  // if neither, null should be returned
  expect(INITIAL_DEPLOYMENT_BLOCK_BY_CHAIN.unset).toBeUndefined();
  expect(await db.getResumeBlockByChain('unset')).toEqual(null);
});
