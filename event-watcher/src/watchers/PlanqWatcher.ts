import { EVMWatcher } from './EVMWatcher';

export class PlanqWatcher extends EVMWatcher {
  constructor() {
    super('planq');
  }

  override async getFinalizedBlockNumber(): Promise<number> {
    const latestBlock = await super.getFinalizedBlockNumber();
    return Math.max(latestBlock - 15, 0);
  }
}
