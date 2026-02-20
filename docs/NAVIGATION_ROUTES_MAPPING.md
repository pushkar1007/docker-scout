# DockerScout Dashboard Navigation & Routes Mapping

This document maps every sidebar navigation item to its corresponding backend route, handler, and permission requirements.

## Overview

**Project**: docker-scout (DockerScout)  
**Frontend Template Engine**: Templ (Go templates)  
**Frontend Sidebar File**: [internal/web/templates/partials/sidebar.templ](../internal/web/templates/partials/sidebar.templ)  
**Backend Routes File**: [internal/web/routes_frontend.go](../internal/web/routes_frontend.go)  
**Handler Package**: [internal/web/handler_*](../internal/web/)

---

## Sidebar Structure & Route Mapping

### 1. OVERVIEW Section
Always expanded, contains top-level infrastructure views.

| Menu Item | Route | Handler | HTTP Method | Required Permission | Feature Flag |
|-----------|-------|---------|-------------|---------------------|--------------|
| **Dashboard** | `/` / `/dashboard` | `DashboardTempl()` | GET | Default (authenticated) | N/A |
| **Infrastructure** | `/overview` | `OverviewTempl()` | GET | Default (authenticated) | N/A |
| **Nodes** | `/nodes` | `HostsTempl()` | GET | `host:view` | N/A |
| **Swarm** | `/swarm` | `SwarmClusterTempl()` | GET | `host:view` | `FeatureSwarm` (EE only) |

**Sidebar Filter**: Swarm hidden by default in Community Edition  
**Active Marker**: `active == "dashboard"`, `active == "overview"`, `active == "nodes"`, `active == "swarm"`

---

### 2. RESOURCES Section
Always expanded, primary Docker resource management.

| Menu Item | Route | Handler | HTTP Method | Required Permission | Notes |
|-----------|-------|---------|-------------|---------------------|-------|
| **Containers** | `/containers` | `ContainersTempl()` | GET | `container:view` | Badge shows running count |
| **Images** | `/images` | `ImagesTempl()` | GET | `image:view` | N/A |
| **Volumes** | `/volumes` | `VolumesTempl()` | GET | `volume:view` | N/A |
| **Networks** | `/networks` | `NetworksTempl()` | GET | `network:view` | N/A |
| **Stacks** | `/stacks` | `StacksTempl()` | GET | `stack:view` | N/A |
| **Templates** | `/container-templates` | `ContainerTemplatesTempl()` | GET | `container:view` | Hidden if `prefs.IsHidden("templates")` |
| **Ports** | `/ports` | `PortsTempl()` | GET | Default (authenticated) | Hidden if `prefs.IsHidden("ports")` |

**Sidebar Filter**: Templates and Ports can be hidden by user preferences  
**Active Marker**: `active == "containers"`, `active == "images"`, etc.

#### Containers Sub-Routes
```
GET     /containers/                          # List all containers
GET     /containers/{id}                      # Detail view
GET     /containers/{id}/stats                # Stats dashboard
GET     /containers/{id}/inspect              # Raw inspect JSON
GET     /containers/{id}/files                # File browser
GET     /containers/{id}/logs                 # Logs view
GET     /containers/{id}/exec                 # Terminal
GET     /containers/{id}/settings             # Settings/config
GET     /containers/new                       # Create form
POST    /containers/create                    # Submit create
POST    /containers/{id}/start                # Action: start
POST    /containers/{id}/stop                 # Action: stop
POST    /containers/{id}/restart              # Action: restart
POST    /containers/{id}/pause                # Action: pause
POST    /containers/{id}/unpause              # Action: unpause
POST    /containers/{id}/kill                 # Action: kill
POST    /containers/{id}/remove               # Action: remove
POST    /containers/{id}/rename               # Action: rename
POST    /containers/bulk/*                    # Bulk operations
```

#### Images Sub-Routes
```
GET     /images/                              # List all images
GET     /images/{id}                          # Detail view
GET     /images/pull                          # Pull form
POST    /images/pull                          # Submit pull
POST    /images/{id}/remove                   # Remove image
POST    /images/prune                         # Prune unused
```

#### Volumes Sub-Routes
```
GET     /volumes/                             # List all volumes
GET     /volumes/{name}                       # Detail/browse
GET     /volumes/new                          # Create form
POST    /volumes/create                       # Submit create
POST    /volumes/{name}/remove                # Remove volume
POST    /volumes/prune                        # Prune unused
```

#### Networks Sub-Routes
```
GET     /networks/                            # List all networks
GET     /networks/{id}                        # Detail view
GET     /networks/new                         # Create form
POST    /networks/create                      # Submit create
POST    /networks/{id}/connect                # Connect container
POST    /networks/{id}/disconnect             # Disconnect container
POST    /networks/{id}/remove                 # Remove network
POST    /networks/prune                       # Prune unused
```

#### Stacks Sub-Routes
```
GET     /stacks/                              # List all stacks
GET     /stacks/catalog                       # Stack catalog
GET     /stacks/catalog/{slug}                # Catalog item detail
GET     /stacks/{name}                        # Stack detail
GET     /stacks/new                           # Create form
POST    /stacks/deploy                        # Deploy from compose
POST    /stacks/catalog/{slug}/deploy         # Deploy from catalog
GET     /stacks/{name}/edit                   # Edit stack
POST    /stacks/{name}/start                  # Start stack
POST    /stacks/{name}/stop                   # Stop stack
POST    /stacks/{name}/restart                # Restart stack
POST    /stacks/{name}/remove                 # Remove stack
```

#### Container Templates Sub-Routes
```
GET     /container-templates/                 # List templates
POST    /container-templates/                 # Create template
POST    /container-templates/{id}/deploy      # Deploy from template
POST    /container-templates/{id}/delete      # Delete template
```

---

### 3. OPERATIONS Section
Collapsible (default expanded), security & operational management.

| Menu Item | Route | Handler | HTTP Method | Required Permission | Badge | Hidden |
|-----------|-------|---------|-------------|---------------------|-------|--------|
| **Security** | `/security` | `SecurityTempl()` | GET | `security:view` | Shows issues count | N/A |
| **Vulnerabilities** | `/vulnerabilities` | `VulnMgmtTempl()` | GET | `security:view` | N/A | N/A |
| **Compliance** | `/compliance` | `ComplianceTempl()` | GET | `security:view` | N/A | `prefs.IsHidden("compliance")` |
| **Updates** | `/updates` | `UpdatesTempl()` | GET | `settings:view` | Shows pending updates | N/A |
| **Backups** | `/backups` | `BackupsTempl()` | GET | `backup:view` | N/A | N/A |
| **Config** | `/config` | `ConfigTempl()` | GET | `config:view` | N/A | `prefs.IsHidden("config")` |
| **Secrets** | `/secrets` | `SecretsTempl()` | GET | Admin only | N/A | `prefs.IsHidden("secrets")` |
| **Lifecycle** | `/lifecycle` | `LifecyclePoliciesTempl()` | GET | `settings:view` | N/A | `prefs.IsHidden("lifecycle")` |
| **Maintenance** | `/maintenance` | `MaintenanceTempl()` | GET | Admin only | N/A | N/A |
| **Bulk Operations** | `/bulk-ops` | `BulkOpsTempl()` | GET | `container:view` | N/A | N/A |

**Sidebar Filter**: Collapsible section controlled by `prefs.IsCollapsed("operations")`

#### Security Sub-Routes
```
GET     /security/                            # Main dashboard (security:view)
GET     /security/trends                      # Trends view
GET     /security/report                      # Report view
GET     /security/container/{id}              # Container security details
POST    /security/scan                        # Scan all (security:scan)
POST    /security/scan/{id}                   # Scan container
POST    /security/issues/{id}/ignore          # Ignore issue
POST    /security/issues/{id}/resolve         # Mark resolved
```

#### Vulnerabilities Sub-Routes
```
GET     /vulnerabilities/                     # Main dashboard (security:view)
POST    /vulnerabilities/scan                 # Scan all (security:scan)
POST    /vulnerabilities/{id}/acknowledge     # Mark acknowledged
POST    /vulnerabilities/{id}/resolve         # Mark resolved
POST    /vulnerabilities/{id}/accept          # Accept risk
```

#### Compliance Sub-Routes
```
GET     /compliance/                          # Main dashboard (security:view)
POST    /compliance/policies                  # Create policy (security:scan)
POST    /compliance/policies/{id}/toggle      # Toggle policy
POST    /compliance/policies/{id}/delete      # Delete policy
POST    /compliance/scan                      # Run scan
POST    /compliance/violations/{id}/acknowledge  # Acknowledge violation
POST    /compliance/violations/{id}/resolve   # Resolve violation
POST    /compliance/violations/{id}/exempt    # Exempt violation
```

#### Updates Sub-Routes
```
GET     /updates/                             # Main dashboard (settings:view)
GET     /updates/{id}/changelog               # Changelog
POST    /updates/check                        # Check for updates (settings:update)
POST    /updates/check-all                    # Check all resources
POST    /updates/{id}/apply                   # Apply single update
POST    /updates/apply-all                    # Apply all updates
POST    /updates/{id}/rollback                # Rollback update
POST    /updates/policies                     # Create auto-update policy
POST    /updates/policies/{id}/toggle         # Toggle policy
POST    /updates/policies/{id}/delete         # Delete policy
```

#### Backups Sub-Routes
```
GET     /backups/                             # List backups (backup:view)
GET     /backups/new                          # Create form
GET     /backups/schedules                    # Schedules list
GET     /backups/{id}                         # Backup detail
GET     /backups/{id}/download                # Download backup
POST    /backups/create                       # Create backup (backup:create)
POST    /backups/{id}/delete                  # Delete backup
POST    /backups/{id}/restore                 # Restore backup (backup:restore)
POST    /backups/schedules                    # Create schedule
POST    /backups/schedules/{id}/delete        # Delete schedule
POST    /backups/schedules/{id}/run           # Run schedule now
```

#### Config Sub-Routes
```
GET     /config/                              # Variables list (config:view)
GET     /config/new                           # Create variable
GET     /config/variables/{id}                # Edit variable
GET     /config/templates                     # Templates list
GET     /config/templates/{id}                # Edit template
GET     /config/audit                         # Audit log
GET     /config/export                        # Export config
POST    /config/variables                     # Create variable (config:create)
POST    /config/variables/{id}                # Update variable (config:update)
POST    /config/variables/{id}/delete         # Delete variable (config:remove)
POST    /config/templates                     # Create template
POST    /config/templates/{id}                # Update template
POST    /config/import                        # Import config
POST    /config/sync/{id}                     # Sync config
```

#### Secrets Sub-Routes
```
GET     /secrets/                             # List secrets (admin only)
POST    /secrets/                             # Create secret
POST    /secrets/{id}/delete                  # Delete secret
POST    /secrets/{id}/rotate                  # Rotate secret
```

#### Lifecycle Sub-Routes
```
GET     /lifecycle/                           # List policies (settings:view)
POST    /lifecycle/policies                   # Create policy (settings:update)
POST    /lifecycle/policies/{id}/toggle       # Toggle policy
POST    /lifecycle/policies/{id}/delete       # Delete policy
POST    /lifecycle/policies/{id}/execute      # Execute policy now
```

#### Maintenance Sub-Routes
```
GET     /maintenance/                         # List windows (admin only)
POST    /maintenance/                         # Create window
POST    /maintenance/{id}/toggle              # Toggle window
POST    /maintenance/{id}/delete              # Delete window
POST    /maintenance/{id}/execute             # Execute now
```

---

### 4. CONNECTIONS Section
Collapsible (default expanded), remote connection management.

| Menu Item | Route | Handler | HTTP Method | Required Permission | Hidden |
|-----------|-------|---------|-------------|---------------------|--------|
| **Databases** | `/connections/database` | `DatabaseConnectionsTempl()` | GET | `host:view` | `prefs.IsHidden("databases")` |
| **LDAP** | `/connections/ldap` | `LDAPConnectionsTempl()` | GET | `host:view` | `prefs.IsHidden("ldap")` |
| **SSH** | `/connections/ssh` | `SSHConnectionsTempl()` | GET | `host:view` | N/A |
| **RDP** | `/connections/rdp` | `RDPConnectionsTempl()` | GET | `host:view` | `prefs.IsHidden("rdp")` |
| **SSH Keys** | `/connections/keys` | `SSHKeysTempl()` | GET | `host:view` | N/A |
| **Shortcuts** | `/connections/shortcuts` | `ShortcutsTempl()` | GET | `host:view` | N/A |

**Sidebar Filter**: Collapsible section controlled by `prefs.IsCollapsed("connections")`

#### Connections Main Route
```
GET     /connections/                        # Main connections dashboard (host:view)
```

#### Database Connections Sub-Routes
```
GET     /connections/database/                # List connections (host:view)
POST    /connections/database/                # Create connection (host:view)
GET     /connections/database/{id}            # Browser view
POST    /connections/database/{id}/test       # Test connection
DELETE  /connections/database/{id}            # Delete connection
POST    /connections/database/{id}/write-mode # Toggle write mode
GET     /connections/database/{id}/query      # Query editor
POST    /connections/database/{id}/query      # Execute query
```

#### LDAP Connections Sub-Routes
```
GET     /connections/ldap/                    # List connections (host:view)
POST    /connections/ldap/                    # Create connection
GET     /connections/ldap/{id}                # Browser view
GET     /connections/ldap/{id}/settings       # Settings
POST    /connections/ldap/{id}/settings       # Update settings
POST    /connections/ldap/{id}/test           # Test connection
DELETE  /connections/ldap/{id}                # Delete connection
POST    /connections/ldap/{id}/write-mode     # Toggle write mode
GET     /connections/ldap/{id}/search         # Search form
POST    /connections/ldap/{id}/search         # Execute search
```

#### SSH Connections Sub-Routes
```
GET     /connections/ssh/                     # List connections (host:view)
GET     /connections/ssh/new                  # Create form
POST    /connections/ssh/                     # Create connection
GET     /connections/ssh/{id}                 # Detail view
POST    /connections/ssh/{id}                 # Update connection
DELETE  /connections/ssh/{id}                 # Delete connection
POST    /connections/ssh/{id}/test            # Test connection
POST    /connections/ssh/{id}/duplicate       # Duplicate connection
GET     /connections/ssh/{id}/terminal        # Terminal session
GET     /connections/ssh/{id}/files           # SFTP browser
GET     /connections/ssh/{id}/files/list      # Browse directory
POST    /connections/ssh/{id}/files/upload    # Upload file
GET     /connections/ssh/{id}/files/download  # Download file
POST    /connections/ssh/{id}/files/delete    # Delete file/folder
POST    /connections/ssh/{id}/files/mkdir     # Create directory
POST    /connections/ssh/{id}/files/rename    # Rename file
GET     /connections/ssh/{id}/tunnels         # SSH tunnels list
POST    /connections/ssh/{id}/tunnels         # Create tunnel
POST    /connections/ssh/{id}/tunnels/{tunnelID}/toggle  # Toggle tunnel
DELETE  /connections/ssh/{id}/tunnels/{tunnelID}         # Delete tunnel
```

#### RDP Connections Sub-Routes
```
GET     /connections/rdp/                     # List connections (host:view)
GET     /connections/rdp/new                  # Create form
POST    /connections/rdp/                     # Create connection
GET     /connections/rdp/{id}                 # Detail view
POST    /connections/rdp/{id}                 # Update connection
DELETE  /connections/rdp/{id}                 # Delete connection
POST    /connections/rdp/{id}/test            # Test connection
GET     /connections/rdp/{id}/download        # Download RDP file
GET     /connections/rdp/{id}/session         # RDP session
```

#### SSH Keys Sub-Routes
```
GET     /connections/keys/                    # List SSH keys (host:view)
GET     /connections/keys/new                 # Create form
POST    /connections/keys/                    # Create key
GET     /connections/keys/{id}                # Edit key
DELETE  /connections/keys/{id}                # Delete key
GET     /connections/keys/{id}/download       # Download key
```

#### Shortcuts Sub-Routes
```
GET     /connections/shortcuts/               # List shortcuts (host:view)
GET     /connections/shortcuts/new            # Create form
POST    /connections/shortcuts/               # Create shortcut
GET     /connections/shortcuts/{id}/edit      # Edit shortcut
POST    /connections/shortcuts/{id}           # Update shortcut
DELETE  /connections/shortcuts/{id}           # Delete shortcut
```

---

### 5. TOOLS Section
Collapsible (default collapsed), productivity & advanced tools.

| Menu Item | Route | Handler | HTTP Method | Required Permission | Hidden |
|-----------|-------|---------|-------------|---------------------|--------|
| **Terminal Hub** | `/terminal` | `TerminalHubTempl()` | GET | `container:exec` | N/A |
| **Editor** | `/editor` | `EditorHub()` (HTMX endpoint) | GET | `host:update` | N/A |
| **Cheat Sheet** | `/tools/cheatsheet` | `CheatSheet()` | GET | Default (authenticated) | N/A |
| **Ansible** | `/tools/ansible` | `AnsibleInventory()` | GET | `host:view` | `prefs.IsHidden("ansible")` |
| **Packet Capture** | `/tools/capture` | `PacketCapture()` | GET | `host:update` | `prefs.IsHidden("capture")` |

**Sidebar Filter**: Collapsible section controlled by `prefs.IsCollapsed("tools")`

#### Terminal Hub Sub-Routes
```
GET     /terminal/                            # Terminal hub (container:exec)
GET     /terminal/picker                      # Container/host picker
WS      /ws/exec/{id}                         # WebSocket exec
WS      /ws/host-exec/{id}                    # Host exec WebSocket
WS      /ws/ssh/{id}                          # SSH WebSocket
GET     /api/v1/terminal/sessions             # Terminal session history
GET     /api/v1/terminal/sessions/active      # Active sessions
GET     /api/v1/terminal/sessions/{id}        # Session details
GET     /api/v1/terminal/sessions/target/{type}/{id}  # Sessions for target
```

#### Editor Sub-Routes
```
GET     /editor/                              # Editor hub (host:update)
GET     /editor/monaco                        # Monaco editor
GET     /editor/nvim                          # Neovim editor
WS      /ws/editor/nvim                       # Neovim WebSocket
GET     /api/snippets/                        # List snippets
POST    /api/snippets/                        # Create snippet
GET     /api/snippets/paths                   # Get snippet paths
GET     /api/snippets/{id}                    # Get snippet
PUT     /api/snippets/{id}                    # Update snippet
DELETE  /api/snippets/{id}                    # Delete snippet
```

#### Cheat Sheet Sub-Routes
```
GET     /tools/cheatsheet                     # Cheat sheet (authenticated)
POST    /tools/cheatsheet/custom              # Create custom command
DELETE  /tools/cheatsheet/custom/{id}         # Delete custom command
```

#### Ansible Inventory Sub-Routes
```
GET     /tools/ansible                        # Inventory browser (host:view)
POST    /tools/ansible/upload                 # Upload inventory
POST    /tools/ansible/parse                  # Parse inventory
DELETE  /tools/ansible/{id}                   # Delete inventory
```

#### Packet Capture Sub-Routes
```
GET     /tools/capture                        # Capture dashboard (host:update)
GET     /tools/capture/{id}                   # Capture detail
POST    /tools/capture/start                  # Start capture
POST    /tools/capture/{id}/stop              # Stop capture
GET     /tools/capture/{id}/download          # Download PCAP
DELETE  /tools/capture/{id}                   # Delete capture
WS      /ws/capture/{id}                      # Capture stream WebSocket
```

---

### 6. INTEGRATIONS Section
Collapsible (default collapsed), git & API integrations.

| Menu Item | Route | Handler | HTTP Method | Required Permission | Hidden |
|-----------|-------|---------|-------------|---------------------|--------|
| **Reverse Proxy** | `/proxy` | `ProxyTempl()` | GET | `host:view` | `prefs.IsHidden("proxy")` |
| **Repositories** | `/integrations/gitea` | `GiteaTempl()` | GET | `stack:view` | `prefs.IsHidden("repositories")` |
| **GitOps** | `/gitops` | `GitOpsTempl()` | GET | `stack:view` | `prefs.IsHidden("gitops")` |
| **Storage** | `/storage` | `StorageTempl()` | GET | `backup:view` | N/A |
| **Notifications** | `/notifications` | `NotificationsTempl()` | GET | Default (authenticated) | N/A |

**Sidebar Filter**: Collapsible section controlled by `prefs.IsCollapsed("integrations")`

#### Reverse Proxy Sub-Routes
```
GET     /proxy/                               # Main dashboard (host:view)
GET     /proxy/setup                          # Setup wizard
POST    /proxy/setup                          # Save setup (host:update)
POST    /proxy/setup/delete                   # Delete setup
POST    /proxy/setup/test                     # Test setup
GET     /proxy/new                            # Create host form
POST    /proxy/hosts                          # Create host (host:update)
POST    /proxy/hosts/{id}                     # Update host
DELETE  /proxy/hosts/{id}                     # Delete host
POST    /proxy/hosts/{id}/enable              # Enable host
POST    /proxy/hosts/{id}/disable             # Disable host
POST    /proxy/sync                           # Sync with proxy
GET     /proxy/{id}                           # Host detail
GET     /proxy/{id}/edit                      # Edit host form

# Certificates sub-route
GET     /proxy/certificates/                  # List certificates
GET     /proxy/certificates/{id}              # Certificate detail
GET     /proxy/certificates/new/letsencrypt   # Let's Encrypt form
GET     /proxy/certificates/new/custom        # Custom cert form
POST    /proxy/certificates/letsencrypt       # Create LE cert (host:update)
POST    /proxy/certificates/custom            # Upload custom cert
POST    /proxy/certificates/{id}/renew        # Renew certificate
DELETE  /proxy/certificates/{id}              # Delete certificate

# Redirections sub-route
GET     /proxy/redirections/                  # List redirections
GET     /proxy/redirections/new               # Create form
POST    /proxy/redirections/                  # Create redirect
GET     /proxy/redirections/{id}/edit         # Edit form
POST    /proxy/redirections/{id}              # Update redirect
DELETE  /proxy/redirections/{id}              # Delete redirect

# Streams sub-route
GET     /proxy/streams/                       # List streams
GET     /proxy/streams/new                    # Create form
POST    /proxy/streams/                       # Create stream
GET     /proxy/streams/{id}/edit              # Edit form
POST    /proxy/streams/{id}                   # Update stream
DELETE  /proxy/streams/{id}                   # Delete stream

# Access Lists sub-route
GET     /proxy/access-lists/                  # List ACLs
GET     /proxy/access-lists/new               # Create form
POST    /proxy/access-lists/                  # Create ACL
GET     /proxy/access-lists/{id}/edit         # Edit form
POST    /proxy/access-lists/{id}              # Update ACL
DELETE  /proxy/access-lists/{id}              # Delete ACL

# Audit sub-route
GET     /proxy/audit                          # Proxy audit log
```

#### Gitea Repositories Sub-Routes
```
GET     /integrations/gitea/                  # Main dashboard (stack:view)
POST    /integrations/gitea/connections       # Create connection
POST    /integrations/gitea/connections/{id}/test     # Test connection
POST    /integrations/gitea/connections/{id}/sync     # Sync repos
POST    /integrations/gitea/connections/{id}/delete   # Delete connection
GET     /integrations/gitea/connections/{id}/templates # Get templates

GET     /integrations/gitea/repos/            # List repos
POST    /integrations/gitea/repos             # Create repo
GET     /integrations/gitea/repos/{id}        # Repo detail
GET     /integrations/gitea/repos/{id}/files  # File browser
GET     /integrations/gitea/repos/{id}/file   # Get file content
POST    /integrations/gitea/repos/{id}/file   # Save file

# Repository management
POST    /integrations/gitea/repos/{id}/edit            # Update repo
POST    /integrations/gitea/repos/{id}/delete          # Delete repo

# Branches, Tags, Commits
GET     /integrations/gitea/repos/{id}/branches       # List branches
POST    /integrations/gitea/repos/{id}/branches       # Create branch
DELETE  /integrations/gitea/repos/{id}/branches/{name} # Delete branch

GET     /integrations/gitea/repos/{id}/tags           # List tags
POST    /integrations/gitea/repos/{id}/tags           # Create tag
DELETE  /integrations/gitea/repos/{id}/tags/{name}     # Delete tag

GET     /integrations/gitea/repos/{id}/commits        # List commits
GET     /integrations/gitea/repos/{id}/commits/{sha}  # Commit detail
GET     /integrations/gitea/repos/{id}/compare        # Compare branches
GET     /integrations/gitea/repos/{id}/diff           # Get diff

# Pull Requests
GET     /integrations/gitea/repos/{id}/pulls/         # List PRs
POST    /integrations/gitea/repos/{id}/pulls/         # Create PR
GET     /integrations/gitea/repos/{id}/pulls/{number} # PR detail
PATCH   /integrations/gitea/repos/{id}/pulls/{number} # Edit PR
POST    /integrations/gitea/repos/{id}/pulls/{number}/merge  # Merge PR

# Issues
GET     /integrations/gitea/repos/{id}/issues/        # List issues
POST    /integrations/gitea/repos/{id}/issues/        # Create issue
GET     /integrations/gitea/repos/{id}/issues/{number} # Issue detail

# More...
GET     /integrations/gitea/repos/{id}/labels         # List labels
GET     /integrations/gitea/repos/{id}/milestones     # List milestones
GET     /integrations/gitea/repos/{id}/collaborators  # List collaborators
```

#### GitOps Sub-Routes
```
GET     /gitops/                              # Main dashboard (stack:view)
POST    /gitops/pipelines                     # Create pipeline (stack:deploy)
POST    /gitops/pipelines/{id}/toggle         # Toggle pipeline
POST    /gitops/pipelines/{id}/delete         # Delete pipeline
POST    /gitops/pipelines/{id}/deploy         # Deploy pipeline
```

#### Storage Sub-Routes
```
GET     /storage/                             # Main dashboard (backup:view)
POST    /storage/connections                  # Create connection (backup:create)
GET     /storage/{connID}/buckets             # List buckets (backup:view)
GET     /storage/{connID}/audit               # Audit log
POST    /storage/{connID}/delete              # Delete connection (backup:create)
POST    /storage/{connID}/test                # Test connection
POST    /storage/{connID}/buckets             # Create bucket

GET     /storage/{connID}/buckets/{bucket}/browse     # Browse bucket
GET     /storage/{connID}/buckets/{bucket}/download   # Download object
GET     /storage/{connID}/buckets/{bucket}/presign-upload  # Presign upload
POST    /storage/{connID}/buckets/{bucket}/delete     # Delete bucket
POST    /storage/{connID}/buckets/{bucket}/upload     # Upload object
POST    /storage/{connID}/buckets/{bucket}/delete-object  # Delete object
POST    /storage/{connID}/buckets/{bucket}/create-folder  # Create folder
```

#### Notifications Sub-Routes
```
GET     /notifications/                       # List notifications (authenticated)
POST    /notifications/mark-all-read         # Mark all read
POST    /notifications/{id}/read             # Mark single read
DELETE  /notifications/{id}                  # Delete notification
```

---

### 7. MONITORING Section
Collapsible (default collapsed), observability & metrics.

| Menu Item | Route | Handler | HTTP Method | Required Permission | Hidden |
|-----------|-------|---------|-------------|---------------------|--------|
| **Metrics** | `/monitoring` | `MonitoringPage()` | GET | `host:view` | `prefs.IsHidden("metrics")` |
| **Alerts** | `/alerts` | `AlertsTempl()` | GET | `security:view` | N/A |
| **Health** | `/health-dashboard` | `HealthDashTempl()` | GET | `container:view` | `prefs.IsHidden("health")` |
| **Logs** | `/logs` | `LogsPageTempl()` | GET | `container:view` | N/A |
| **Log Analysis** | `/logs/management` | `LogManagement()` | GET | `container:view` | `prefs.IsHidden("log-management")` |
| **Events** | `/events` | `EventsTempl()` | GET | Default (authenticated) | N/A |
| **Topology** | `/topology` | `TopologyTempl()` | GET | Default (authenticated) | N/A |
| **Dependencies** | `/dependencies` | `DependenciesTempl()` | GET | Default (authenticated) | `prefs.IsHidden("dependencies")` |

**Sidebar Filter**: Collapsible section controlled by `prefs.IsCollapsed("monitoring")`

#### Monitoring Sub-Routes
```
GET     /monitoring/                          # Prometheus/Grafana dashboard (host:view)
GET     /monitoring/{id}                      # Container-specific metrics
GET     /partials/monitoring/host             # Host metrics partial
GET     /partials/monitoring/containers       # Container metrics partial
GET     /partials/monitoring/history          # Metrics history JSON
WS      /ws/monitoring/stats                  # Metrics stream WebSocket
WS      /ws/monitoring/container/{id}         # Container metrics stream
```

#### Alerts Sub-Routes
```
GET     /alerts/                              # List alerts (security:view)
GET     /alerts/{id}                          # Edit alert form
POST    /alerts/                              # Create alert (security:scan)
POST    /alerts/{id}                          # Update alert
DELETE  /alerts/{id}                          # Delete alert
POST    /alerts/{id}/enable                   # Enable alert
POST    /alerts/{id}/disable                  # Disable alert
POST    /alerts/events/{id}/ack               # Acknowledge alert event
POST    /alerts/silences                      # Create silence
DELETE  /alerts/silences/{id}                 # Delete silence
```

#### Logs Sub-Routes
```
GET     /logs/                                # Centralized logs dashboard (container:view)
GET     /logs/management/                     # Log management interface
POST    /logs/uploads/                        # Upload log file
GET     /logs/uploads/{id}                    # Analyze uploaded log
DELETE  /logs/uploads/{id}                    # Delete uploaded log
GET     /api/logs/search                      # Search logs API
```

#### Events Sub-Routes
```
GET     /events/                              # Events stream (authenticated)
WS      /ws/events                            # Events WebSocket
```

#### Topology Sub-Routes
```
GET     /topology/                            # Topology visualization (authenticated)
```

#### Dependencies Sub-Routes
```
GET     /dependencies/                        # Dependency graph (authenticated)
```

---

### 8. ADMIN Section
Only visible to admin users, always expanded.

| Menu Item | Route | Handler | HTTP Method | Required Permission | Sub-Routes |
|-----------|-------|---------|-------------|---------------------|------------|
| **Teams** | `/teams` | `TeamsTempl()` | GET | Admin + `user:view` | Yes |
| **Users** | `/users` | `UsersTempl()` | GET | Admin + `user:view` | Yes |
| **Roles** | `/admin/roles` | `RolesTempl()` | GET | Admin + `role:view` | Yes |
| **OAuth** | `/admin/oauth` | `OAuthProvidersTempl()` | GET | Admin + FeatureOAuth | Yes |
| **LDAP** | `/admin/ldap` | `LDAPProvidersTempl()` | GET | Admin | Yes |
| **Channels** | `/admin/notifications/channels` | `NotificationChannelsTempl()` | GET | Admin | Yes |
| **Quotas** | `/quotas` | `QuotasTempl()` | GET | Admin | Yes |
| **Access Audit** | `/access-audit` | `AccessAuditTempl()` | GET | Admin | Yes |
| **Settings** | `/settings` | `SettingsTempl()` | GET | Admin + `settings:view` | Yes |
| **Webhooks** | `/webhooks` | `WebhooksTempl()` | GET | Admin | Yes |
| **Runbooks** | `/runbooks` | `RunbooksTempl()` | GET | Admin | Yes |
| **Registries** | `/registries` | `RegistriesTempl()` | GET | Admin | Yes |
| **Jobs** | `/jobs` | `JobsTempl()` | GET | Admin | Yes |
| **License** | `/license` | `LicenseTempl()` | GET | Admin | Yes |

**Visibility**: Only shown when `user.Role == "admin"`

#### Teams Sub-Routes
```
GET     /teams/                               # List teams (admin + user:view)
GET     /teams/new                            # Create form
POST    /teams/                               # Create team (admin + user:create)
GET     /teams/{id}                           # Team detail
GET     /teams/{id}/edit                      # Edit form
POST    /teams/{id}                           # Update team
DELETE  /teams/{id}                           # Delete team
POST    /teams/{id}/members                   # Add member
DELETE  /teams/{id}/members/{userID}          # Remove member
POST    /teams/{id}/permissions               # Grant permission
DELETE  /teams/{id}/permissions/{permID}      # Revoke permission
```

#### Users Sub-Routes
```
GET     /users/                               # List users (admin + user:view)
GET     /users/new                            # Create form
POST    /users/                               # Create user (admin + user:create)
GET     /users/{id}                           # Edit user form
POST    /users/{id}                           # Update user (admin + user:update)
POST    /users/{id}/enable                    # Enable user
POST    /users/{id}/disable                   # Disable user
DELETE  /users/{id}                           # Delete user (admin + user:remove)
```

#### Roles Sub-Routes
```
GET     /admin/roles/                         # List roles (admin + role:view)
POST    /admin/roles/                         # Create role (admin + role:create)
GET     /admin/roles/{id}                     # Edit role form
POST    /admin/roles/{id}                     # Update role (admin + role:update)
POST    /admin/roles/{id}/enable              # Enable role
POST    /admin/roles/{id}/disable             # Disable role
DELETE  /admin/roles/{id}                     # Delete role (admin + role:remove)
```

#### OAuth Providers Sub-Routes *(Enterprise only)*
```
GET     /admin/oauth/                         # List providers (admin, EE only)
POST    /admin/oauth/                         # Create provider (admin)
GET     /admin/oauth/{id}                     # Edit form
POST    /admin/oauth/{id}                     # Update provider
DELETE  /admin/oauth/{id}                     # Delete provider
POST    /admin/oauth/{id}/enable              # Enable provider
POST    /admin/oauth/{id}/disable             # Disable provider
```

#### LDAP Providers Sub-Routes
```
GET     /admin/ldap/                          # List providers (admin)
POST    /admin/ldap/                          # Create provider (admin)
GET     /admin/ldap/{id}                      # Edit form
POST    /admin/ldap/{id}                      # Update provider
DELETE  /admin/ldap/{id}                      # Delete provider
POST    /admin/ldap/{id}/enable               # Enable provider
POST    /admin/ldap/{id}/disable              # Disable provider
POST    /admin/ldap/{id}/test                 # Test connection
```

#### Notification Channels Sub-Routes
```
GET     /admin/notifications/channels/        # List channels (admin)
POST    /admin/notifications/channels/        # Create channel (admin)
GET     /admin/notifications/channels/{name}/edit  # Edit form
POST    /admin/notifications/channels/{name}/delete # Delete channel
POST    /admin/notifications/channels/{name}/test  # Test channel
```

#### Quotas Sub-Routes
```
GET     /quotas/                              # List quotas (admin)
POST    /quotas/                              # Create quota (admin)
POST    /quotas/{id}/toggle                   # Enable/disable quota
POST    /quotas/{id}/delete                   # Delete quota
```

#### Access Audit Sub-Routes
```
GET     /access-audit/                        # List audit records (admin)
POST    /access-audit/export                  # Export audit log
POST    /access-audit/sessions/{id}/revoke    # Revoke session
```

#### Settings Sub-Routes
```
GET     /settings/                            # General settings (admin + settings:view)
GET     /settings/totp                        # TOTP setup page
POST    /settings/                            # Update settings (admin + settings:update)
POST    /settings/totp/verify                 # Verify TOTP setup
POST    /settings/totp/disable                # Disable TOTP
```

#### Webhooks Sub-Routes
```
GET     /webhooks/                            # List webhooks (admin)
POST    /webhooks/                            # Create webhook (admin)
POST    /webhooks/{id}/delete                 # Delete webhook
POST    /webhooks/autodeploy                  # Create auto-deploy webhook
POST    /webhooks/autodeploy/{id}/delete      # Delete auto-deploy
```

#### Runbooks Sub-Routes
```
GET     /runbooks/                            # List runbooks (admin)
POST    /runbooks/                            # Create runbook (admin)
POST    /runbooks/{id}/delete                 # Delete runbook
POST    /runbooks/{id}/execute                # Execute runbook
```

#### Registries Sub-Routes
```
GET     /registries/                          # List registries (admin)
POST    /registries/                          # Create registry (admin)
POST    /registries/{id}/update               # Update registry
POST    /registries/{id}/delete               # Delete registry
```

#### Jobs Sub-Routes
```
GET     /jobs/                                # List jobs (admin)
GET     /jobs/{id}                            # Job detail
POST    /jobs/{id}/cancel                     # Cancel job
POST    /jobs/{id}/delete                     # Delete job
WS      /ws/jobs/{id}                         # Job progress WebSocket
```

#### License Sub-Routes
```
GET     /license/                             # License info (admin)
POST    /license/activate                     # Activate license
POST    /license/deactivate                   # Deactivate license
```

---

## Profile Routes

Available to all authenticated users, not shown in sidebar.

```
GET     /profile/                         # Profile page (authenticated)
POST    /profile/                         # Update profile
POST    /profile/password                 # Change password
PUT     /profile/preferences              # Update preferences
PUT     /profile/sidebar-prefs            # Update sidebar preferences (AJAX)
POST    /profile/preferences/reset        # Reset preferences
POST    /profile/theme                    # Toggle dark/light theme
POST    /profile/export                   # Export user data
DELETE  /profile/                         # Delete account
DELETE  /profile/sessions                 # Delete all sessions
DELETE  /profile/sessions/{id}            # Delete specific session
GET     /profile/language                 # Get language preference
POST    /profile/language                 # Set language
```

---

## WebSocket Routes

Real-time data streaming endpoints (all require authentication + permission).

```
WS      /ws/logs/{id}                     # Container logs stream (container:logs)
WS      /ws/exec/{id}                     # Container exec/shell (container:exec)
WS      /ws/stats/{id}                    # Container stats (container:view)
WS      /ws/host-exec/{id}                # Host shell (host:update)
WS      /ws/events                        # System events (container:view)
WS      /ws/jobs/{id}                     # Job progress (container:view)
WS      /ws/capture/{id}                  # Packet capture stream (host:update)
WS      /ws/metrics                       # Metrics stream (host:view)
WS      /ws/monitoring/stats              # Monitoring stats (host:view)
WS      /ws/monitoring/container/{id}     # Container metrics (host:view)
WS      /ws/editor/nvim                   # Neovim editor (host:update)
WS      /ws/ssh/{id}                      # SSH terminal (host:view)
WS      /ws/rdp/{id}                      # RDP session (host:view)
```

---

## API Routes

REST API endpoints (mostly JSON responses).

```
# Snippets API (editor storage)
GET     /api/snippets/                    # List snippets (host:update)
POST    /api/snippets/                    # Create snippet
GET     /api/snippets/{id}                # Get snippet
PUT     /api/snippets/{id}                # Update snippet
DELETE  /api/snippets/{id}                # Delete snippet

# Logs API
GET     /api/logs/search                  # Search logs

# Host API
GET     /api/v1/hosts/{hostID}/browse     # Browse filesystem (host:update)
GET     /api/v1/hosts/{hostID}/file/*     # Read file
GET     /api/v1/hosts/{hostID}/download/* # Download file
POST    /api/v1/hosts/{hostID}/mkdir/*    # Create directory
DELETE  /api/v1/hosts/{hostID}/file/*     # Delete file

# Terminal API
GET     /api/v1/terminal/sessions         # List terminal sessions
GET     /api/v1/terminal/sessions/active  # Active sessions
GET     /api/v1/terminal/sessions/{id}    # Session details
GET     /api/v1/terminal/sessions/target/{type}/{id}  # Sessions for target

# Compliance API (EE)
GET     /api/v1/compliance/frameworks     # List frameworks (FeatureCompliance)
POST    /api/v1/compliance/frameworks/{id}/assess  # Run assessment
GET     /api/v1/compliance/assessments/{assessmentId}/report

# OPA API (EE)
GET     /api/v1/opa/policies              # Get OPA policies (FeatureOPAPolicies)
POST    /api/v1/opa/evaluate/container/{id}  # Evaluate container

# Log Aggregation API (EE)
GET     /api/v1/logs/search               # Search aggregated logs (FeatureLogAggregation)
GET     /api/v1/logs/stats                # Log statistics

# Image Signing API (EE)
GET     /api/v1/images/signing/signatures # List image signatures (FeatureImageSigning)
GET     /api/v1/images/signing/verify     # Verify image
GET     /api/v1/images/signing/trust-policies  # Trust policies

# Runtime Security API (EE)
GET     /api/v1/runtime-security/events   # Runtime events (FeatureRuntimeSecurity)
GET     /api/v1/runtime-security/dashboard # Runtime dashboard
GET     /api/v1/runtime-security/rules    # Runtime rules

# Git Sync API (EE)
GET     /api/v1/git-sync/configs          # List git sync configs (FeatureGitSync)
GET     /api/v1/git-sync/stats            # Git sync statistics

# Ephemeral Environments API (EE)
GET     /api/v1/ephemeral/environments    # List ephemeral envs (FeatureEphemeralEnvs)
POST    /api/v1/ephemeral/environments    # Create ephemeral env

# Manifest Builder API (EE)
GET     /api/v1/manifests/templates       # List manifest templates (FeatureManifestBuilder)
POST    /api/v1/manifests/generate        # Generate manifest
POST    /api/v1/manifests/validate        # Validate manifest
```

---

## HTMX Partials

Fragment/snippet endpoints for dynamic updates.

```
GET     /partials/stats                        # Stats cards partial
GET     /partials/containers                   # Containers list partial
GET     /partials/container/{id}               # Single container row
GET     /partials/images                       # Images list partial
GET     /partials/events                       # Events list partial
GET     /partials/notifications                # Notifications partial
GET     /partials/search                       # Search results partial

# Monitoring partials
GET     /partials/monitoring/host              # Host metrics
GET     /partials/monitoring/containers        # Container metrics
GET     /partials/monitoring/history           # Metrics history (JSON)
```

---

## Permission Matrix

| Permission | Scope | Implied Actions |
|-----------|-------|-----------------|
| `container:view` | List/view containers | GET /containers, /containers/{id}, /containers/{id}/stats |
| `container:create` | Create/update containers | POST /containers/create, /containers/{id}/settings |
| `container:start` | Control container state | POST /containers/{id}/start/stop/restart/pause |
| `container:stop` | Stop/pause containers | POST /containers/{id}/stop/pause/kill |
| `container:remove` | Delete containers | DELETE /containers/{id} |
| `container:logs` | View logs | WS /ws/logs/{id}, GET /containers/{id}/logs |
| `container:exec` | Execute commands | GET /terminal, WS /ws/exec/{id} |
| `image:view` | List/view images | GET /images, /images/{id} |
| `image:pull` | Pull images | POST /images/pull |
| `image:remove` | Delete images | DELETE /images/{id} |
| `volume:view` | List/view volumes | GET /volumes, /volumes/{name} |
| `volume:create` | Create volumes | POST /volumes/create |
| `volume:remove` | Delete volumes | DELETE /volumes/{name} |
| `network:view` | List/view networks | GET /networks, /networks/{id} |
| `network:create` | Create/connect networks | POST /networks/create, /networks/{id}/connect |
| `network:remove` | Delete networks | DELETE /networks/{id} |
| `stack:view` | List/view stacks | GET /stacks, /integrations/gitea, /gitops |
| `stack:deploy` | Deploy stacks | POST /stacks/deploy, /gitops/pipelines |
| `stack:update` | Edit stacks | POST /stacks/{name}/edit |
| `stack:remove` | Delete stacks | DELETE /stacks/{name} |
| `host:view` | List/view hosts/nodes | GET /nodes, /overview, /connections/* |
| `host:create` | Create agents/proxies | POST /nodes/create, /proxy/* |
| `host:update` | Edit hosts/proxies | POST /nodes/{id}/edit, PUT /proxy/* |
| `host:remove` | Delete hosts | DELETE /nodes/{id} |
| `security:view` | View security dashboards | GET /security, /vulnerabilities, /compliance |
| `security:scan` | Run scans | POST /security/scan, /vulnerabilities/scan |
| `backup:view` | List backups | GET /backups |
| `backup:create` | Create backups | POST /backups/create, /storage/* |
| `backup:restore` | Restore backups | POST /backups/{id}/restore |
| `config:view` | View configs | GET /config, /config/audit |
| `config:create` | Create variables | POST /config/variables |
| `config:update` | Update variables | POST /config/variables/{id} |
| `config:remove` | Delete variables | DELETE /config/variables/{id} |
| `settings:view` | View settings | GET /settings, /updates |
| `settings:update` | Modify settings | POST /settings, /updates/apply |
| `user:view` | List users | GET /users (admin) |
| `user:create` | Create users | POST /users (admin) |
| `user:update` | Edit users | POST /users/{id} (admin) |
| `user:remove` | Delete users | DELETE /users/{id} (admin) |
| `role:view` | List roles | GET /admin/roles (admin) |
| `role:create` | Create roles | POST /admin/roles (admin) |
| `role:update` | Edit roles | POST /admin/roles/{id} (admin) |
| `role:remove` | Delete roles | DELETE /admin/roles/{id} (admin) |

---

## Feature Flags (License-Based)

| Feature | Route | Requirement | Notes |
|---------|-------|-------------|-------|
| `FeatureSwarm` | `/swarm/*` | Enterprise | Swarm mode & clustering |
| `FeatureOAuth` | `/admin/oauth/*` | Enterprise | OAuth 2.0 providers |
| `FeatureOPAPolicies` | `/api/v1/opa/*` | Enterprise | OPA policy engine |
| `FeatureCompliance` | `/api/v1/compliance/*` | Enterprise | Compliance frameworks |
| `FeatureLogAggregation` | `/api/v1/logs/*` | Enterprise | Log aggregation |
| `FeatureImageSigning` | `/api/v1/images/signing/*` | Enterprise | Image signing & verification |
| `FeatureRuntimeSecurity` | `/api/v1/runtime-security/*` | Enterprise | Runtime security monitoring |
| `FeatureCustomDashboards` | `/api/v1/dashboards/*` | Enterprise | Custom dashboard layouts |
| `FeatureGitSync` | `/api/v1/git-sync/*` | Enterprise | GitOps bidirectional sync |
| `FeatureEphemeralEnvs` | `/api/v1/ephemeral/*` | Enterprise | Branch-based ephemeral environments |
| `FeatureManifestBuilder` | `/api/v1/manifests/*` | Enterprise | Visual GitOps manifest builder |

---

## Active Parameter Usage

The `active` parameter in the sidebar determines which menu item is highlighted. Values map to route paths:

```
active = "dashboard"           when on /, /dashboard
active = "overview"            when on /overview
active = "nodes"               when on /nodes
active = "swarm"               when on /swarm
active = "containers"          when on /containers
active = "images"              when on /images
active = "volumes"             when on /volumes
active = "networks"            when on /networks
active = "stacks"              when on /stacks
active = "container-templates" when on /container-templates
active = "ports"               when on /ports
active = "security"            when on /security
active = "vulnerabilities"     when on /vulnerabilities
active = "compliance"          when on /compliance
active = "updates"             when on /updates
active = "backups"             when on /backups
active = "config"              when on /config
active = "secrets"             when on /secrets
active = "lifecycle"           when on /lifecycle
active = "maintenance"         when on /maintenance
active = "bulk-ops"            when on /bulk-ops
active = "connections-database"  when on /connections/database
active = "connections-ldap"    when on /connections/ldap
active = "connections-ssh"     when on /connections/ssh
active = "connections-rdp"     when on /connections/rdp
active = "connections-keys"    when on /connections/keys
active = "connections-shortcuts" when on /connections/shortcuts
active = "terminal"            when on /terminal
active = "editor"              when on /editor
active = "cheatsheet"          when on /tools/cheatsheet
active = "ansible"             when on /tools/ansible
active = "capture"             when on /tools/capture
active = "proxy"               when on /proxy
active = "gitea"               when on /integrations/gitea
active = "gitops"              when on /gitops
active = "storage"             when on /storage
active = "notifications"       when on /notifications
active = "monitoring"          when on /monitoring
active = "alerts"              when on /alerts
active = "health"              when on /health-dashboard
active = "logs"                when on /logs
active = "log-management"      when on /logs/management
active = "events"              when on /events
active = "topology"            when on /topology
active = "dependencies"        when on /dependencies
active = "teams"               when on /teams (admin)
active = "users"               when on /users (admin)
active = "roles"               when on /admin/roles (admin)
active = "oauth-providers"     when on /admin/oauth (admin)
active = "ldap-providers"      when on /admin/ldap (admin)
active = "notification-channels" when on /admin/notifications/channels (admin)
active = "quotas"              when on /quotas (admin)
active = "access-audit"        when on /access-audit (admin)
active = "settings"            when on /settings (admin)
active = "webhooks"            when on /webhooks (admin)
active = "runbooks"            when on /runbooks (admin)
active = "registries"          when on /registries (admin)
active = "jobs"                when on /jobs (admin)
active = "license"             when on /license (admin)
```

---

## Location of Sidebar Preferences

**File**: [internal/web/templates/types/types.go](../internal/web/templates/types/types.go)

```go
type SidebarPreferences struct {
    CollapsedSections map[string]bool  // Sections: operations, connections, tools, integrations, monitoring
    HiddenItems       map[string]bool  // Individual items: swarm, templates, ports, etc.
}
```

**Collapse Defaults** (controlled by functions like `IsCollapsed(section string)`):
- `operations` - Default expanded
- `connections` - Default expanded
- `tools` - Default collapsed
- `integrations` - Default collapsed
- `monitoring` - Default collapsed

**Hidden Item Defaults** (can be toggled via `prefs.IsHidden(item string)`):
- `swarm` - True in CE, False in EE
- `templates` - False
- `ports` - False
- `compliance` - False
- `config` - False
- `secrets` - False
- `lifecycle` - False
- `maintenance` - False
- `ansible` - False
- `capture` - False
- `proxy` - False
- `repositories` - False
- `gitops` - False
- `metrics` - False
- `health` - False
- `log-management` - False
- `dependencies` - False
- `oauth` - False
- `ldap-providers` - False
- `channels` - False
- `quotas` - False
- `webhooks` - False
- `runbooks` - False
- `registries` - False
- `rdp` - False

---

## Navigation Badge System

Badges are displayed on menu items showing counts (e.g., security issues, running containers).

**Function**: [sidebar.templ](../internal/web/templates/partials/sidebar.templ#L380-L395) - `getBadge(stats *types.StatsData, key string) int`

```go
case "containers":
    return stats.ContainersRunning
case "security":
    return stats.SecurityIssuesCount
case "updates":
    return stats.UpdatesAvailableCount
```

**Badge Styles**:
- `navItem()` - Gray badge (neutral)
- `navItemDanger()` - Red badge (security issues)
- `navItemPrimary()` - Primary colored badge (updates)

---

## Related Documentation

- [Architecture Overview](./architecture.md)
- [API Documentation](./api.md)
- [Development Guide](./development.md)
- [Security Model](./security.md) *(if exists)*

---

## Version Info

- **Last Updated**: Current Session
- **Framework**: Go Chi Router + Templ
- **License**: AGPLv3 (Community), Commercial (Business/Enterprise)
- **Edition Support**: CE, Business, Enterprise
