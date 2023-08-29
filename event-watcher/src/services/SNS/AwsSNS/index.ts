import crypto from 'node:crypto';
import {
  SNSClient,
  PublishCommand,
  PublishCommandInput,
  PublishBatchCommand,
  PublishBatchCommandInput,
  PublishBatchRequestEntry,
} from '@aws-sdk/client-sns';
import { AwsSNSConfig, SNSInput, SNSPublishMessageOutput } from '../types';
import BaseSNS from '../BaseSNS';
import { env } from '../../../config';

const isDev = env.NODE_ENV === 'development';

class AwsSNS extends BaseSNS {
  private client: SNSClient;
  private subject: string;
  private topicArn: string;

  constructor(private config: AwsSNSConfig) {
    super();

    const { region, credentials, subject, topicArn } = this.config;

    this.subject = subject;
    this.topicArn = topicArn;
    const credentialsConfig = {
      region,
      ...(isDev && { credentials }),
    };

    this.client = new SNSClient(credentialsConfig);
  }

  override async publishMessage(
    { message, subject, groupId, deduplicationId }: SNSInput,
    fifo: boolean = false,
  ): Promise<SNSPublishMessageOutput> {
    const input: PublishCommandInput = {
      TopicArn: this.topicArn!,
      Subject: subject ?? this.subject!,
      Message: message,
      ...(fifo && { MessageGroupId: groupId }),
      ...(fifo && { MessageDeduplicationId: deduplicationId }),
    };

    try {
      const command = new PublishCommand(input);
      await this.client?.send(command);
    } catch (error) {
      console.error(error);

      return {
        status: 'error',
      };
    }

    return {
      status: 'success',
    };
  }

  override async publishMessages(
    messages: SNSInput[],
    fifo: boolean = false,
  ): Promise<SNSPublishMessageOutput> {
    const CHUNK_SIZE = 10;
    const batches: PublishBatchCommandInput[] = [];
    const inputs: PublishBatchRequestEntry[] = messages.map(
      ({ message, subject, groupId, deduplicationId }) => ({
        Id: crypto.randomUUID(),
        Subject: subject ?? this.subject!,
        Message: message,
        ...(fifo && { MessageGroupId: groupId }),
        ...(fifo && { MessageDeduplicationId: deduplicationId }),
      }),
    );

    // PublishBatchCommand: only supports max 10 items per batch
    for (let i = 0; i <= inputs.length; i += CHUNK_SIZE) {
      const batch: PublishBatchCommandInput = {
        TopicArn: this.topicArn!,
        PublishBatchRequestEntries: inputs.slice(i, i + CHUNK_SIZE),
      };

      batches.push(batch);
    }

    try {
      const promises = [];
      const errors = [];
      for (const batch of batches) {
        const command = new PublishBatchCommand(batch);
        promises.push(this.client?.send(command));
      }

      const results = await Promise.allSettled(promises);

      for (const result of results) {
        if (result.status !== 'fulfilled') {
          console.error(result.reason);
          errors.push(result.reason);
        }
      }

      if (errors.length > 0) {
        console.error(errors);

        return {
          status: 'error',
          reasons: errors,
        };
      }
    } catch (error) {
      console.error(error);

      return {
        status: 'error',
      };
    }

    return {
      status: 'success',
    };
  }
}

export default AwsSNS;
