import { Router } from "express";
import {PlService} from "./pl.service";

export class PlController {
  private router: Router = Router();
  private readonly plService: PlService = new PlService();

  constructor(prefix: string) {
    this.initializeRoutes(prefix);
  }

  private initializeRoutes(prefix: string): void {
    this.router.get(`${prefix}`, (req, res) => this.plService.getStatus(req, res));
  }

  public getRouter(): Router {
    return this.router;
  }
}
