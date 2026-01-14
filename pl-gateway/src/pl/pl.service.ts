import { DirectSecp256k1HdWallet, DirectSecp256k1Wallet, Registry } from "@cosmjs/proto-signing";
import { createProtobufRpcClient, GasPrice, ProtobufRpcClient, QueryClient, SigningStargateClient, StargateClient } from "@cosmjs/stargate";
import { Tendermint34Client } from "@cosmjs/tendermint-rpc";
import { Request, Response } from "express";
import { Entity } from "src/proto/pl/factory/v1/entity";
import { MsgCreateEntity } from "src/proto/pl/factory/v1/tx";
import { createBufRegistry } from "src/utils/message";
import * as factory from "src/proto/pl/factory/v1/query";


export class PlService {
  private reigstry: Registry = createBufRegistry(
    [
      // Tx messages
      ['/pl.factory.v1.MsgCreateEntity', MsgCreateEntity],
      // Types
      ['/pl.factory.v1.Entity', Entity],
    ] as any
  )

  private async getQueryclient(): Promise<QueryClient> {
    const tmClient = await Tendermint34Client.connect(process.env.RPC_ENDPOINT || "http://localhost:26657");
    return new QueryClient(tmClient);
  }

  private async getProtobufRpcClient(): Promise<ProtobufRpcClient> {
      return createProtobufRpcClient(await this.getQueryclient());
  }

  private async getStargateClient(): Promise<StargateClient> {
      return await StargateClient.connect(process.env.RPC_ENDPOINT || "http://localhost:26657");
  }


  public getStatus(req: Request, res: Response): Response<any, Record<string, any>> {
    return res.status(200).send({
      status: true,
      internetProtocol: req.ip
    });
  }

  public async listToken(req: factory.QueryAllEntityRequest) {
    try {
      const client = await this.getProtobufRpcClient();

      const queryCurrency = new factory.QueryClientImpl(client)

      return {
        success: true,
        data: {
          storage: await queryCurrency.ListEntity(req)
        }
      };
    } catch (e) {
      return {
        success: false,
        error: e
      };
    }
}

public async createEntity(wallet: DirectSecp256k1HdWallet | DirectSecp256k1Wallet, hash: string) {
    if (!wallet || (!hash || hash.length == 0)) {
      return {
        success: false,
        error: "wallet, denom, amount, to is required"
      }
    }
    
    try {
      const [w] = await wallet.getAccounts()

      const client = await SigningStargateClient.connectWithSigner(
        process.env.RPC_ENDPOINT,
        wallet,
        { gasPrice: GasPrice.fromString(process.env.GAS_FEE), registry: this.reigstry }
      );

      const createEntityMsg = {
        typeUrl: '/example.currency.v1.MsgMintToken',
        value: {
          creator: w.address,
          hash,
        }
      };
      
      const result = await client.signAndBroadcast(
          w.address,
          [createEntityMsg],
          'auto',
          'Create Entity'
      );

      client.disconnect();
      return {success: true, result};
    } catch (e) {
      return {
        success: false,
        error: e
      }
    }
  }

}
