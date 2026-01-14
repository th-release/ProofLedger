import { Router } from "express";
import {PlService} from "./pl.service";
import { WalletService } from "src/wallet/wallet.service";

export class PlController {
  private router: Router = Router();
  private readonly plService: PlService = new PlService();
  private readonly walletService: WalletService = new WalletService()

  constructor(prefix: string) {
    this.initializeRoutes(prefix);
  }

  private initializeRoutes(prefix: string): void {
    this.router.get(`${prefix}`, (req, res) => this.plService.getStatus(req, res));
    
    this.router.get(`${prefix}/list`, async(req, res) => {
      const { offset = "0", limit = "10" } = req.query

      const result = await this.plService.listEntity({
        pagination: {
          key: new Uint8Array(0),
          offset: +offset,
          limit: +limit,
          countTotal: true,
          reverse: false,
        }
      })

      return res.status(200).json(result)
    })

    this.router.get(`${prefix}/detail/:clid`, async(req, res) => {
      const result = await this.plService.detailEntity({
        clid: req.params.clid
      })

      return res.status(200).json(result)
    })

    this.router.post(`${prefix}/mnemonic/create`, async(req, res) => {
      const { mnemonic, hash } = req.body

      const wallet = await this.walletService.walletFromMnemonic(mnemonic)

      if (!wallet.success) {
        return res.status(500).json(wallet)
      }

      const result = await this.plService.createEntity(wallet.wallet, hash)

      if (!result.success) {
        return res.status(500).json(result)
      }

      return res.status(201).json({
        success: result.success,
        response: result.result.msgResponses,
        gasUsed: result.result.gasUsed.toString(),
        gasWanted: result.result.gasWanted.toString(),
        height: result.result.height.toString(),
        txHash: result.result.transactionHash,
        txIndex: result.result.txIndex,
      })
    })

    this.router.post(`${prefix}/privatekey/create`, async(req, res) => {
      const { privatekey, hash } = req.body

      const wallet = await this.walletService.walletFromPrivateKey(privatekey)

      if (!wallet.success) {
        return res.status(500).json(wallet)
      }

      const result = await this.plService.createEntity(wallet.wallet, hash)

      if (!result.success) {
        return res.status(500).json(result)
      }

      return res.status(201).json({
        success: result.success,
        response: result.result.msgResponses,
        gasUsed: result.result.gasUsed.toString(),
        gasWanted: result.result.gasWanted.toString(),
        height: result.result.height.toString(),
        txHash: result.result.transactionHash,
        txIndex: result.result.txIndex,
      })
    })
  }

  public getRouter(): Router {
    return this.router;
  }
}
