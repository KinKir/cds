package workflow

import (
	"database/sql"
	"time"

	"github.com/ovh/cds/engine/api/database/gorpmapping"

	"github.com/ovh/cds/sdk"
)

type Workflows []Workflow

func (ws Workflows) Get() sdk.Workflows {
	res := make(sdk.Workflows, len(ws))
	for i, w := range ws {
		res[i] = w.Get()
	}
	return res
}

// Workflow is a gorp wrapper around sdk.WorkflowData
type Workflow struct {
	sdk.Workflow
	ProjectKey string `db:"project_key"`
}

// Notification is a gorp wrapper around sdk.WorkflowNotification
type Notification sdk.WorkflowNotification

// Run is a gorp wrapper around sdk.WorkflowRun
type Run sdk.WorkflowRun

// Coverage is a gorp wrapper around sdk.WorkflowNodeRunCoverage
type Coverage sdk.WorkflowNodeRunCoverage

type dbNodeRunVulenrabilitiesReport sdk.WorkflowNodeRunVulnerabilityReport

// NodeRun is a gorp wrapper around sdk.WorkflowNodeRun
type NodeRun struct {
	WorkflowID             sql.NullInt64  `db:"workflow_id"`
	WorkflowRunID          int64          `db:"workflow_run_id"`
	ApplicationID          sql.NullInt64  `db:"application_id"`
	ID                     int64          `db:"id"`
	WorkflowNodeID         int64          `db:"workflow_node_id"`
	WorkflowNodeName       string         `db:"workflow_node_name"`
	Number                 int64          `db:"num"`
	SubNumber              int64          `db:"sub_num"`
	Status                 string         `db:"status"`
	Start                  time.Time      `db:"start"`
	Done                   time.Time      `db:"done"`
	LastModified           time.Time      `db:"last_modified"`
	HookEvent              sql.NullString `db:"hook_event"`
	Manual                 sql.NullString `db:"manual"`
	SourceNodeRuns         sql.NullString `db:"source_node_runs"`
	Payload                sql.NullString `db:"payload"`
	PipelineParameters     sql.NullString `db:"pipeline_parameters"`
	BuildParameters        sql.NullString `db:"build_parameters"`
	Tests                  sql.NullString `db:"tests"`
	Commits                sql.NullString `db:"commits"`
	Stages                 sql.NullString `db:"stages"`
	TriggersRun            sql.NullString `db:"triggers_run"`
	VCSRepository          sql.NullString `db:"vcs_repository"`
	VCSBranch              sql.NullString `db:"vcs_branch"`
	VCSTag                 sql.NullString `db:"vcs_tag"`
	VCSHash                sql.NullString `db:"vcs_hash"`
	VCSServer              sql.NullString `db:"vcs_server"`
	Header                 sql.NullString `db:"header"`
	UUID                   sql.NullString `db:"uuid"`
	OutgoingHook           sql.NullString `db:"outgoinghook"`
	HookExecutionTimestamp sql.NullInt64  `db:"hook_execution_timestamp"`
	ExecutionID            sql.NullString `db:"execution_id"`
	Callback               sql.NullString `db:"callback"`
}

// JobRun is a gorp wrapper around sdk.WorkflowNodeJobRun
type JobRun struct {
	ProjectID                 int64          `db:"project_id"`
	ID                        int64          `db:"id"`
	WorkflowNodeRunID         int64          `db:"workflow_node_run_id"`
	Job                       sql.NullString `db:"job"`
	Parameters                sql.NullString `db:"variables"`
	Status                    string         `db:"status"`
	Retry                     int            `db:"retry"`
	Queued                    time.Time      `db:"queued"`
	Start                     time.Time      `db:"start"`
	Done                      time.Time      `db:"done"`
	Model                     string         `db:"model"`
	ExecGroups                sql.NullString `db:"exec_groups"`
	IntegrationPluginBinaries sql.NullString `db:"integration_plugin_binaries"`
	BookedBy                  sdk.Service    `db:"-"`
	ContainsService           bool           `db:"contains_service"`
	ModelType                 sql.NullString `db:"model_type"`
	Header                    sql.NullString `db:"header"`
	HatcheryName              string         `db:"hatchery_name"`
	WorkerName                string         `db:"worker_name"`
}

// ToJobRun transform the JobRun with data of the provided sdk.WorkflowNodeJobRun
func (j *JobRun) ToJobRun(jr *sdk.WorkflowNodeJobRun) (err error) {
	j.ProjectID = jr.ProjectID
	j.ID = jr.ID
	j.WorkflowNodeRunID = jr.WorkflowNodeRunID
	j.Job, err = gorpmapping.JSONToNullString(jr.Job)
	if err != nil {
		return sdk.WrapError(err, "column job")
	}
	j.Parameters, err = gorpmapping.JSONToNullString(jr.Parameters)
	if err != nil {
		return sdk.WrapError(err, "column variables")
	}
	j.Status = jr.Status
	j.Retry = jr.Retry
	j.Queued = jr.Queued
	j.Start = jr.Start
	j.Done = jr.Done
	j.Model = jr.Model
	j.ModelType = sql.NullString{Valid: true, String: string(jr.ModelType)}
	j.ContainsService = jr.ContainsService
	j.ExecGroups, err = gorpmapping.JSONToNullString(jr.ExecGroups)
	j.WorkerName = jr.WorkerName
	j.HatcheryName = jr.HatcheryName
	if err != nil {
		return sdk.WrapError(err, "column exec_groups")
	}
	j.IntegrationPluginBinaries, err = gorpmapping.JSONToNullString(jr.IntegrationPluginBinaries)
	if err != nil {
		return sdk.WrapError(err, "column integration_plugin_binaries")
	}
	j.Header, err = gorpmapping.JSONToNullString(jr.Header)
	if err != nil {
		return sdk.WrapError(err, "column header")
	}
	return nil
}

// WorkflowNodeRunJob returns a sdk.WorkflowNodeRunJob
func (j JobRun) WorkflowNodeRunJob() (sdk.WorkflowNodeJobRun, error) {
	jr := sdk.WorkflowNodeJobRun{
		ProjectID:         j.ProjectID,
		ID:                j.ID,
		WorkflowNodeRunID: j.WorkflowNodeRunID,
		Status:            j.Status,
		Retry:             j.Retry,
		Queued:            j.Queued,
		QueuedSeconds:     time.Now().Unix() - j.Queued.Unix(),
		Start:             j.Start,
		Done:              j.Done,
		BookedBy:          j.BookedBy,
		ContainsService:   j.ContainsService,
		HatcheryName:      j.HatcheryName,
		WorkerName:        j.WorkerName,
		Model:             j.Model,
	}
	if err := gorpmapping.JSONNullString(j.Job, &jr.Job); err != nil {
		return jr, sdk.WrapError(err, "column job")
	}
	if err := gorpmapping.JSONNullString(j.Parameters, &jr.Parameters); err != nil {
		return jr, sdk.WrapError(err, "column variables")
	}
	if err := gorpmapping.JSONNullString(j.ExecGroups, &jr.ExecGroups); err != nil {
		return jr, sdk.WrapError(err, "column exec_groups")
	}
	if err := gorpmapping.JSONNullString(j.IntegrationPluginBinaries, &jr.IntegrationPluginBinaries); err != nil {
		return jr, sdk.WrapError(err, "integration_plugin_binaries")
	}
	if err := gorpmapping.JSONNullString(j.Header, &jr.Header); err != nil {
		return jr, sdk.WrapError(err, "header")
	}
	if j.ModelType.Valid {
		jr.ModelType = j.ModelType.String
	}
	if defaultOS != "" && defaultArch != "" {
		var modelFound, osArchFound bool
		for _, req := range jr.Job.Action.Requirements {
			if req.Type == sdk.ModelRequirement {
				modelFound = true
			}
			if req.Type == sdk.OSArchRequirement {
				osArchFound = true
			}
		}

		if !modelFound && !osArchFound {
			jr.Job.Action.Requirements = append(jr.Job.Action.Requirements, sdk.Requirement{
				Name:  defaultOS + "/" + defaultArch,
				Type:  sdk.OSArchRequirement,
				Value: defaultOS + "/" + defaultArch,
			})
		}
	}
	if defaultRegion != "" {
		var regionFound bool
		for _, req := range jr.Job.Action.Requirements {
			if req.Type == sdk.RegionRequirement {
				regionFound = true
				break
			}
		}

		if !regionFound {
			jr.Job.Action.Requirements = append(jr.Job.Action.Requirements, sdk.Requirement{
				Name:  defaultRegion,
				Type:  sdk.RegionRequirement,
				Value: defaultRegion,
			})
		}
	}
	return jr, nil
}

// NodeRunArtifact is a gorp wrapper around sdk.WorkflowNodeRunArtifact
type NodeRunArtifact sdk.WorkflowNodeRunArtifact

// dbStaticFiles is a gorp wrapper around sdk.StaticFiles
type dbStaticFiles sdk.StaticFiles

// RunTag is a gorp wrapper around sdk.WorkflowRunTag
type RunTag sdk.WorkflowRunTag

// hookModel is a gorp wrapper around sdk.WorkflowHookModel
type hookModel sdk.WorkflowHookModel

// outgoingHookModel is a gorp wrapper around sdk.WorkflowHookModel
type outgoingHookModel sdk.WorkflowHookModel

type auditWorkflow sdk.AuditWorkflow

type dbNodeData sdk.Node
type dbNodeContextData sqlNodeContextData
type dbNodeTriggerData sdk.NodeTrigger
type dbNodeOutGoingHookData sdk.NodeOutGoingHook
type dbNodeJoinData sdk.NodeJoin
type dbNodeHookData sdk.NodeHook

type dbAsCodeEvents sdk.AsCodeEvent

func init() {
	gorpmapping.Register(gorpmapping.New(Workflow{}, "workflow", true, "id"))
	gorpmapping.Register(gorpmapping.New(Run{}, "workflow_run", true, "id"))
	gorpmapping.Register(gorpmapping.New(NodeRun{}, "workflow_node_run", true, "id"))
	gorpmapping.Register(gorpmapping.New(JobRun{}, "workflow_node_run_job", true, "id"))
	gorpmapping.Register(gorpmapping.New(NodeRunArtifact{}, "workflow_node_run_artifacts", true, "id"))
	gorpmapping.Register(gorpmapping.New(RunTag{}, "workflow_run_tag", false, "workflow_run_id", "tag"))
	gorpmapping.Register(gorpmapping.New(hookModel{}, "workflow_hook_model", true, "id"))
	gorpmapping.Register(gorpmapping.New(outgoingHookModel{}, "workflow_outgoing_hook_model", true, "id"))
	gorpmapping.Register(gorpmapping.New(Notification{}, "workflow_notification", true, "id"))
	gorpmapping.Register(gorpmapping.New(auditWorkflow{}, "workflow_audit", true, "id"))
	gorpmapping.Register(gorpmapping.New(Coverage{}, "workflow_node_run_coverage", false, "workflow_id", "workflow_run_id", "workflow_node_run_id", "repository", "branch"))
	gorpmapping.Register(gorpmapping.New(dbStaticFiles{}, "workflow_node_run_static_files", true, "id"))
	gorpmapping.Register(gorpmapping.New(dbNodeRunVulenrabilitiesReport{}, "workflow_node_run_vulnerability", true, "id"))
	gorpmapping.Register(gorpmapping.New(dbNodeData{}, "w_node", true, "id"))
	gorpmapping.Register(gorpmapping.New(dbNodeHookData{}, "w_node_hook", true, "id"))
	gorpmapping.Register(gorpmapping.New(dbNodeContextData{}, "w_node_context", true, "id"))
	gorpmapping.Register(gorpmapping.New(dbNodeTriggerData{}, "w_node_trigger", true, "id"))
	gorpmapping.Register(gorpmapping.New(dbNodeOutGoingHookData{}, "w_node_outgoing_hook", true, "id"))
	gorpmapping.Register(gorpmapping.New(dbNodeJoinData{}, "w_node_join", true, "id"))
	gorpmapping.Register(gorpmapping.New(dbAsCodeEvents{}, "as_code_events", true, "id"))
}
