import { ApiService } from './api-service.js';
import { Sequence } from '../models/sequence.js';
import { Remediation } from '../models/remediation.js';
import { SequenceTypes } from '../models/sequence-types.js';
import { Trace } from '../models/trace.js';
import { EventTypes } from '../models/event-types.js';
import { DeploymentInformation, Service } from '../models/service.js';
import { Approval } from '../models/approval.js';
import { Project } from '../models/project.js';
import { ResultTypes } from '../models/result-types.js';
import { EventState } from '../models/event-state.js';

export class DataService {
  private apiService;
  private readonly MAX_SEQUENCE_PAGE_SIZE = 100;
  private readonly MAX_TRACE_PAGE_SIZE = 50;

  constructor(apiUrl: string, apiToken: string) {
    this.apiService = new ApiService(apiUrl, apiToken);
  }

  public async getProject(projectName: string, includeRemediation: boolean, includeApproval: boolean): Promise<Project> {
    const response = await this.apiService.getProject(projectName);
    const project = Project.fromJSON(response.data);
    let remediations: Remediation[] = [];

    if (includeRemediation) {
      remediations = await this.getRemediations(projectName);
    }
    for (const stage of project.stages) {
      for (const service of stage.services) {
        const keptnContext = service.getLatestSequence(stage.stageName);
        if (keptnContext) {
          try {
            await this.fetchServiceDetails(service, stage.stageName, keptnContext, projectName, includeRemediation, includeApproval, remediations);
          }
          catch (error) {
            console.error(error);
          }
        }
      }
    }
    return project;
  }

  private async fetchServiceDetails(service: Service, stageName: string, keptnContext: string, projectName: string, includeRemediation: boolean, includeApproval: boolean, remediations: Remediation[]): Promise<void> {
    service.latestSequence = await this.getSequence(projectName, stageName, keptnContext, true);
    service.deploymentInformation = await this.getDeploymentInformation(service, projectName, stageName);
    if (includeRemediation) {
      const serviceRemediations = remediations.filter(remediation => remediation.service === service.serviceName && remediation.stages.some(s => s.name === stageName));
      for (const remediation of serviceRemediations) {
        remediation.reduceToStage(stageName);
      }
      service.openRemediations = serviceRemediations;
    }
    if (includeApproval) {
      service.openApprovals = await this.getApprovals(projectName, stageName, service.serviceName);
    }
  }

  public async getSequence(projectName: string, stageName?: string, keptnContext?: string, includeEvaluationTrace = false): Promise<Sequence | undefined> {
    const response = await this.apiService.getSequences(projectName, 1, undefined, undefined, undefined, undefined, keptnContext);
    let sequence = response.data.states[0];
    if (sequence && stageName) { // we just need the result of a stage
      sequence = Sequence.fromJSON(sequence);
      sequence.reduceToStage(stageName);
      if (includeEvaluationTrace) {
        const stage = sequence.stages.find(s => s.name === stageName);
        if (stage) { // get latest evaluation
          const evaluationTraces = await this.getEvaluationResults(projectName, sequence.service, stageName, 1, sequence.shkeptncontext);
          if (evaluationTraces) {
            stage.latestEvaluationTrace = evaluationTraces.shift();
          }
        }
      }
    }
    return Sequence.fromJSON(sequence);
  }

  public async getDeploymentInformation(service: Service, projectName: string, stageName: string): Promise<DeploymentInformation | undefined> {
    const result = await this.apiService.getTracesWithResult(EventTypes.DEPLOYMENT_FINISHED, 1, projectName, stageName, service.serviceName, ResultTypes.PASSED);
    const traceData = result.data.events[0];
    let deploymentInformation: DeploymentInformation | undefined;
    if (traceData) {
      const trace = Trace.fromJSON(traceData);
      deploymentInformation = {
        deploymentUrl: trace.getDeploymentUrl(),
        image: trace.getShortImageName()
      };
    }
    return deploymentInformation;
  }

  private async getEvaluationResults(projectName: string, serviceName: string, stageName: string, pageSize: number, keptnContext?: string): Promise<Trace[]> {
    const response = await this.apiService.getEvaluationResults(projectName, serviceName, stageName, pageSize, keptnContext);
    return response.data.events.map(trace => Trace.fromJSON(trace));
  }

  public async getSequences(projectName: string, sequenceName: string, stageName?: string, keptnContext?: string): Promise<Sequence[]> {
    const response = await this.apiService.getSequences(projectName, this.MAX_SEQUENCE_PAGE_SIZE, sequenceName, undefined, undefined, undefined, keptnContext);
    const sequences = response.data.states;
    for (const sequence of sequences) {
      if (stageName) { // we just need the result of a stage
        if (sequence.name === SequenceTypes.REMEDIATION) { // if the sequence is a remediation also return the problemTitle
          const traceResponse = await this.apiService.getTraces(this.buildRemediationEvent(stageName), 1, projectName, stageName, sequence.service);
          const trace = traceResponse.data.events[0];
          sequence.problemTitle = trace?.data.problem?.ProblemTitle;
        }
        sequence.stages = sequence.stages.filter(stage => stage.name === stageName);
      }
    }
    return sequences.map(sequence => Sequence.fromJSON(sequence));
  }

  public async getRemediations(projectName: string): Promise<Remediation[]> {
    const sequences = await this.getSequences(projectName, SequenceTypes.REMEDIATION);
    const remediations: Remediation[] = [];
    for (const sequence of sequences) {
      const stageName = sequence.stages[0].name;
      const response = await this.apiService.getTraces(this.buildRemediationEvent(stageName), this.MAX_TRACE_PAGE_SIZE, projectName, stageName, sequence.service);
      const traces = response.data.events;
      const stage = {...sequence.stages[0], actions: []};
      const remediation: Remediation = Remediation.fromJSON({...sequence, stages: [stage]});

      remediation.problemTitle = traces[0]?.data.problem?.ProblemTitle;
      for (const trace of traces) {
        if (trace.type === EventTypes.ACTION_TRIGGERED && trace.data.action) {
          const finishedAction = traces.find(t => t.triggeredid === trace.id && t.type === EventTypes.ACTION_FINISHED);
          const startedAction = traces.find(t => t.triggeredid === trace.id && t.type === EventTypes.ACTION_STARTED);
          let state: EventState;
          if (finishedAction) {
            state = EventState.FINISHED;
          }
          else if (startedAction) {
            state = EventState.STARTED;
          }
          else {
            state = EventState.TRIGGERED;
          }

          remediation.stages[0].actions.push({...trace.data.action, state, result: finishedAction?.data.result});
        }
      }
      remediations.push(remediation);
    }
    return remediations;
  }

  private async getTrace(keptnContext: string, projectName: string, stageName: string, serviceName: string, eventType: EventTypes): Promise<Trace | undefined> {
    const response = await this.apiService.getTraces(eventType, 1, projectName, stageName, serviceName, keptnContext);
    return response.data.events.shift();
  }

  public async getApprovals(projectName: string, stageName: string, serviceName: string): Promise<Approval[]> {
    // currently a workaround to fetch open approvals
    const responseTriggered = await this.apiService.getTraces(EventTypes.APPROVAL_TRIGGERED, this.MAX_TRACE_PAGE_SIZE, projectName, stageName, serviceName);
    const responseStarted = await this.apiService.getTraces(EventTypes.APPROVAL_STARTED, this.MAX_TRACE_PAGE_SIZE, projectName, stageName, serviceName);
    const tracesTriggered = responseTriggered.data.events.filter(triggeredTrace => !responseStarted.data.events.some(startedTrace => startedTrace.triggeredid === triggeredTrace.id));
    const approvals: Approval[] = [];
    // for each approval the latest evaluation trace (before this event) is needed
    for (const trace of tracesTriggered) {
      const evaluationTrace = await this.getTrace(trace.shkeptncontext, projectName, stageName, serviceName, EventTypes.EVALUATION_FINISHED);
      approvals.push({
        evaluationTrace,
        trace
      });
    }
    return approvals;
  }

  private buildRemediationEvent(stageName: string) {
    return `sh.keptn.event.${stageName}.${SequenceTypes.REMEDIATION}.${EventState.TRIGGERED}`;
  }
}
