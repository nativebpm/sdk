package nativebpm

import (
	"context"
	"os"

	"gitlab.com/nativebpm/platform/sdk/go/api"
)

// Client wraps the generated APIClient to provide a Fluent API.
type Client struct {
	apiClient *api.APIClient
}

// NewClient creates a new Client instance configured with host and API token.
func NewClient(hostURL, apiToken string) (*Client, error) {
	cfg := api.NewConfiguration()
	cfg.Servers = api.ServerConfigurations{
		{
			URL: hostURL,
		},
	}
	if apiToken != "" {
		cfg.AddDefaultHeader("Authorization", "Bearer "+apiToken)
	}
	return &Client{
		apiClient: api.NewAPIClient(cfg),
	}, nil
}

// Definitions returns the Process Definitions management service.
func (c *Client) Definitions() *DefinitionsService {
	return &DefinitionsService{client: c}
}

// Instances returns the Process Instances management service.
func (c *Client) Instances() *InstancesService {
	return &InstancesService{client: c}
}

// Tasks returns the Human Tasks management service.
func (c *Client) Tasks() *TasksService {
	return &TasksService{client: c}
}

// Incidents returns the Incidents management service.
func (c *Client) Incidents() *IncidentsService {
	return &IncidentsService{client: c}
}

// Webhooks returns the Outgoing Webhooks management service.
func (c *Client) Webhooks() *WebhooksService {
	return &WebhooksService{client: c}
}

// DefinitionsService manages process schemas
type DefinitionsService struct {
	client *Client
}

type ListDefinitionsBuilder struct {
	service *DefinitionsService
}

func (s *DefinitionsService) List() *ListDefinitionsBuilder {
	return &ListDefinitionsBuilder{service: s}
}

func (b *ListDefinitionsBuilder) Send(ctx context.Context) ([]api.ProcessDefinition, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.ListDefinitions(ctx).Execute()
	return res, err
}

type DeployDefinitionBuilder struct {
	service *DefinitionsService
	bpmnXML []byte
	err     error
}

func (s *DefinitionsService) Deploy() *DeployDefinitionBuilder {
	return &DeployDefinitionBuilder{service: s}
}

func (b *DeployDefinitionBuilder) WithID(id string) *DeployDefinitionBuilder {
	return b
}

func (b *DeployDefinitionBuilder) WithName(name string) *DeployDefinitionBuilder {
	return b
}

func (b *DeployDefinitionBuilder) WithBPMN(xmlData []byte) *DeployDefinitionBuilder {
	b.bpmnXML = xmlData
	return b
}

func (b *DeployDefinitionBuilder) Send(ctx context.Context) (*api.ProcessDefinition, error) {
	if b.err != nil {
		return nil, b.err
	}
	tmpFile, err := os.CreateTemp("", "*.bpmn")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write(b.bpmnXML); err != nil {
		return nil, err
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, err
	}

	res, _, err := b.service.client.apiClient.DefaultAPI.DeployDefinition(ctx).
		File(tmpFile).
		Execute()
	return res, err
}

// InstancesService manages process runs
type InstancesService struct {
	client *Client
}

func (s *InstancesService) Incidents() *IncidentsService {
	return &IncidentsService{client: s.client}
}

type ListInstancesBuilder struct {
	service *InstancesService
}

func (s *InstancesService) List() *ListInstancesBuilder {
	return &ListInstancesBuilder{service: s}
}

func (b *ListInstancesBuilder) Send(ctx context.Context) ([]api.ProcessInstance, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.ListInstances(ctx).Execute()
	return res, err
}

type StartInstanceBuilder struct {
	service      *InstancesService
	definitionID string
	instanceID   string
	businessKey  string
	variables    map[string]interface{}
}

func (s *InstancesService) Start(definitionID string) *StartInstanceBuilder {
	return &StartInstanceBuilder{
		service:      s,
		definitionID: definitionID,
	}
}

func (b *StartInstanceBuilder) WithInstanceID(id string) *StartInstanceBuilder {
	b.instanceID = id
	return b
}

func (b *StartInstanceBuilder) WithBusinessKey(key string) *StartInstanceBuilder {
	b.businessKey = key
	return b
}

func (b *StartInstanceBuilder) WithVariable(key string, val interface{}) *StartInstanceBuilder {
	if b.variables == nil {
		b.variables = make(map[string]interface{})
	}
	b.variables[key] = val
	return b
}

func (b *StartInstanceBuilder) Send(ctx context.Context) (*api.ProcessInstance, error) {
	req := api.StartInstanceRequest{
		InstanceId:  &b.instanceID,
		BusinessKey: &b.businessKey,
		Variables:   b.variables,
	}
	res, _, err := b.service.client.apiClient.DefaultAPI.StartInstance(ctx, b.definitionID).
		StartInstanceRequest(req).
		Execute()
	return res, err
}

type GetInstanceBuilder struct {
	service *InstancesService
	id      string
}

func (s *InstancesService) Get(id string) *GetInstanceBuilder {
	return &GetInstanceBuilder{service: s, id: id}
}

func (b *GetInstanceBuilder) Send(ctx context.Context) (*api.ProcessInstance, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.GetInstance(ctx, b.id).Execute()
	return res, err
}

type GetInstanceHistoryBuilder struct {
	service *InstancesService
	id      string
}

func (s *InstancesService) History(id string) *GetInstanceHistoryBuilder {
	return &GetInstanceHistoryBuilder{service: s, id: id}
}

func (b *GetInstanceHistoryBuilder) Send(ctx context.Context) ([]api.HistoryRecord, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.GetInstanceHistory(ctx, b.id).Execute()
	return res, err
}

type CompleteInstanceTaskBuilder struct {
	service   *InstancesService
	id        string
	nodeID    string
	variables map[string]interface{}
}

func (s *InstancesService) CompleteTask(id string) *CompleteInstanceTaskBuilder {
	return &CompleteInstanceTaskBuilder{service: s, id: id}
}

func (b *CompleteInstanceTaskBuilder) WithNodeID(nodeID string) *CompleteInstanceTaskBuilder {
	b.nodeID = nodeID
	return b
}

func (b *CompleteInstanceTaskBuilder) WithVariable(key string, val interface{}) *CompleteInstanceTaskBuilder {
	if b.variables == nil {
		b.variables = make(map[string]interface{})
	}
	b.variables[key] = val
	return b
}

func (b *CompleteInstanceTaskBuilder) Send(ctx context.Context) (*api.ProcessInstance, error) {
	req := api.CompleteInstanceTaskRequest{
		NodeId:    b.nodeID,
		Variables: b.variables,
	}
	res, _, err := b.service.client.apiClient.DefaultAPI.CompleteInstanceTask(ctx, b.id).
		CompleteInstanceTaskRequest(req).
		Execute()
	return res, err
}

type ResumeInstanceBuilder struct {
	service *InstancesService
	id      string
}

func (s *InstancesService) Resume(id string) *ResumeInstanceBuilder {
	return &ResumeInstanceBuilder{service: s, id: id}
}

func (b *ResumeInstanceBuilder) Send(ctx context.Context) (*api.ProcessInstance, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.ResumeInstance(ctx, b.id).Execute()
	return res, err
}

// TasksService manages human steps
type TasksService struct {
	client *Client
}

type ListTasksBuilder struct {
	service  *TasksService
	status   string
	assignee string
}

func (s *TasksService) List() *ListTasksBuilder {
	return &ListTasksBuilder{service: s}
}

func (b *ListTasksBuilder) WithStatus(status string) *ListTasksBuilder {
	b.status = status
	return b
}

func (b *ListTasksBuilder) WithAssignee(assignee string) *ListTasksBuilder {
	b.assignee = assignee
	return b
}

func (b *ListTasksBuilder) Send(ctx context.Context) ([]api.TaskRecord, error) {
	req := b.service.client.apiClient.DefaultAPI.ListTasks(ctx)
	if b.status != "" {
		req = req.Status(b.status)
	}
	if b.assignee != "" {
		req = req.Assignee(b.assignee)
	}
	res, _, err := req.Execute()
	return res, err
}

type ClaimTaskBuilder struct {
	service  *TasksService
	id       string
	assignee string
}

func (s *TasksService) Claim(id string) *ClaimTaskBuilder {
	return &ClaimTaskBuilder{service: s, id: id}
}

func (b *ClaimTaskBuilder) WithAssignee(assignee string) *ClaimTaskBuilder {
	b.assignee = assignee
	return b
}

func (b *ClaimTaskBuilder) Send(ctx context.Context) (*api.TaskRecord, error) {
	req := api.ClaimTaskRequest{
		Assignee: b.assignee,
	}
	res, _, err := b.service.client.apiClient.DefaultAPI.ClaimTask(ctx, b.id).
		ClaimTaskRequest(req).
		Execute()
	return res, err
}

type CompleteTaskBuilder struct {
	service   *TasksService
	id        string
	variables map[string]interface{}
}

func (s *TasksService) Complete(id string) *CompleteTaskBuilder {
	return &CompleteTaskBuilder{service: s, id: id}
}

func (b *CompleteTaskBuilder) WithVariable(key string, val interface{}) *CompleteTaskBuilder {
	if b.variables == nil {
		b.variables = make(map[string]interface{})
	}
	b.variables[key] = val
	return b
}

func (b *CompleteTaskBuilder) Send(ctx context.Context) (*api.ProcessInstance, error) {
	req := api.CompleteTaskRequest{
		Variables: b.variables,
	}
	res, _, err := b.service.client.apiClient.DefaultAPI.CompleteTask(ctx, b.id).
		CompleteTaskRequest(req).
		Execute()
	return res, err
}

// IncidentsService manages active incidents
type IncidentsService struct {
	client *Client
}

type ListIncidentsBuilder struct {
	service           *IncidentsService
	processInstanceID string
}

func (s *IncidentsService) List(processInstanceID string) *ListIncidentsBuilder {
	return &ListIncidentsBuilder{service: s, processInstanceID: processInstanceID}
}

func (b *ListIncidentsBuilder) Send(ctx context.Context) ([]api.IncidentRecord, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.ListIncidents(ctx, b.processInstanceID).Execute()
	return res, err
}

type ResolveIncidentBuilder struct {
	service           *IncidentsService
	processInstanceID string
	incidentID        string
}

func (s *IncidentsService) Resolve(processInstanceID, incidentID string) *ResolveIncidentBuilder {
	return &ResolveIncidentBuilder{
		service:           s,
		processInstanceID: processInstanceID,
		incidentID:        incidentID,
	}
}

func (b *ResolveIncidentBuilder) Send(ctx context.Context) (string, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.ResolveIncident(ctx, b.processInstanceID, b.incidentID).Execute()
	if err != nil {
		return "", err
	}
	return res.GetStatus(), nil
}

// WebhooksService manages outbound events integration
type WebhooksService struct {
	client *Client
}

type ListWebhooksBuilder struct {
	service *WebhooksService
}

func (s *WebhooksService) List() *ListWebhooksBuilder {
	return &ListWebhooksBuilder{service: s}
}

func (b *ListWebhooksBuilder) Send(ctx context.Context) ([]api.WebhookRecord, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.ListWebhooks(ctx).Execute()
	return res, err
}

type CreateWebhookBuilder struct {
	service     *WebhooksService
	targetURL   string
	secretToken string
	eventTypes  []string
}

func (s *WebhooksService) Create() *CreateWebhookBuilder {
	return &CreateWebhookBuilder{service: s}
}

func (b *CreateWebhookBuilder) WithURL(targetURL string) *CreateWebhookBuilder {
	b.targetURL = targetURL
	return b
}

func (b *CreateWebhookBuilder) WithEvents(events []string) *CreateWebhookBuilder {
	b.eventTypes = events
	return b
}

func (b *CreateWebhookBuilder) Send(ctx context.Context) (*api.WebhookRecord, error) {
	req := api.CreateWebhookRequest{
		Url:    b.targetURL,
		Secret: &b.secretToken,
		Events: b.eventTypes,
	}
	res, _, err := b.service.client.apiClient.DefaultAPI.CreateWebhook(ctx).
		CreateWebhookRequest(req).
		Execute()
	return res, err
}

type UpdateWebhookBuilder struct {
	service     *WebhooksService
	id          string
	targetURL   string
	secretToken string
	eventTypes  []string
}

func (s *WebhooksService) Update(id string) *UpdateWebhookBuilder {
	return &UpdateWebhookBuilder{service: s, id: id}
}

func (b *UpdateWebhookBuilder) WithURL(targetURL string) *UpdateWebhookBuilder {
	b.targetURL = targetURL
	return b
}

func (b *UpdateWebhookBuilder) Send(ctx context.Context) (*api.WebhookRecord, error) {
	req := api.CreateWebhookRequest{
		Url:    b.targetURL,
		Secret: &b.secretToken,
		Events: b.eventTypes,
	}
	res, _, err := b.service.client.apiClient.DefaultAPI.UpdateWebhook(ctx, b.id).
		CreateWebhookRequest(req).
		Execute()
	return res, err
}

type DeleteWebhookBuilder struct {
	service   *WebhooksService
	webhookID string
}

func (s *WebhooksService) Delete(id string) *DeleteWebhookBuilder {
	return &DeleteWebhookBuilder{service: s, webhookID: id}
}

func (b *DeleteWebhookBuilder) Send(ctx context.Context) (string, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.DeleteWebhook(ctx, b.webhookID).Execute()
	if err != nil {
		return "", err
	}
	return res.GetStatus(), nil
}

type ListWebhookDeliveriesBuilder struct {
	service   *WebhooksService
	webhookID string
}

func (s *WebhooksService) Deliveries(webhookID string) *ListWebhookDeliveriesBuilder {
	return &ListWebhookDeliveriesBuilder{service: s, webhookID: webhookID}
}

func (b *ListWebhookDeliveriesBuilder) Send(ctx context.Context) ([]api.WebhookDeliveryRecord, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.ListWebhookDeliveries(ctx, b.webhookID).Execute()
	return res, err
}

type TestWebhookBuilder struct {
	service   *WebhooksService
	webhookID string
}

func (s *WebhooksService) Test(webhookID string) *TestWebhookBuilder {
	return &TestWebhookBuilder{service: s, webhookID: webhookID}
}

func (b *TestWebhookBuilder) Send(ctx context.Context) (*api.ResolveIncident200Response, error) {
	res, _, err := b.service.client.apiClient.DefaultAPI.TestWebhook(ctx, b.webhookID).Execute()
	return res, err
}
