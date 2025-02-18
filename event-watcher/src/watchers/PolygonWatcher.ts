import { EVMWatcher } from './EVMWatcher';

export class PolygonWatcher extends EVMWatcher {
  constructor() {
    super('polygon');
  }

  // override async getFinalizedBlockNumber(): Promise<number> {
  //   this.logger.debug('fetching last child block from Ethereum');
  //   const rootChain = new ethers.utils.Interface([
  //     `function getLastChildBlock() external view returns (uint256)`,
  //   ]);
  //   const callData = rootChain.encodeFunctionData('getLastChildBlock');
  //   const callResult = (
  //     await this.http.post(
  //       POLYGON_ROOT_CHAIN_RPC,
  //       [
  //         {
  //           jsonrpc: '2.0',
  //           id: 1,
  //           method: 'eth_call',
  //           params: [
  //             { to: POLYGON_ROOT_CHAIN_ADDRESS, data: callData },
  //             'latest', // does the phylax use latest?
  //           ],
  //         },
  //       ],
  //       AXIOS_CONFIG_JSON,
  //     )
  //   )?.data?.[0]?.result;
  //   const block = rootChain.decodeFunctionResult('getLastChildBlock', callResult)[0].toNumber();
  //   this.logger.debug(`rooted child block ${block}`);
  //   return block;
  // }
}
