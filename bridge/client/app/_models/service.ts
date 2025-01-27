import {Root} from './root';
import {Trace} from './trace';
import { Deployment } from './deployment';
import {EventTypes} from './event-types';
import {Sequence} from './sequence';

export class Service {
  serviceName!: string;
  deployedImage?: string;
  stage!: string;
  allDeploymentsLoaded = false;
  deployments: Deployment[] = [];
  lastEventTypes?: {[key: string]: {eventId: string, keptnContext: string, time: number}};
  sequences: Sequence[] = [];
  roots: Root[] = [];
  openApprovals: Trace[] = [];

  static fromJSON(data: unknown) {
    return Object.assign(new this(), data);
  }

  get deploymentContext(): string | undefined {
    return this.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.keptnContext ?? this.evaluationContext;
  }

  get deploymentTime(): number | undefined {
    return this.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.time || this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED]?.time;
  }

  get evaluationContext(): string | undefined {
    return this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED]?.keptnContext;
  }
  public getShortImageName(): string | undefined {
    return this.deployedImage?.split('/').pop()?.split(':').find(() => true);
  }

  getImageName(): string | undefined {
    return this.deployedImage?.split('/').pop();
  }

  getImageVersion(): string | undefined {
    return this.deployedImage?.split(':').pop();
  }

  getOpenApprovals(): Trace[] {
    return this.openApprovals || [];
  }

  getOpenProblems(): Root[] {
    // show running remediation or last faulty remediation
    return this.roots?.filter((root, index) => root.isRemediation() && (!root.isFinished() || root.isFaulty() && index === 0)) || [];
  }

  getRecentRoot(): Root {
    return this.roots[0];
  }

  getRecentEvaluation(): Trace | undefined {
    return this.getRecentRoot()?.getEvaluation(this.stage);
  }

  public hasRemediations(): boolean {
    return this.deployments.some(d => d.stages.some(s => s.remediations.length !== 0));
  }
}
