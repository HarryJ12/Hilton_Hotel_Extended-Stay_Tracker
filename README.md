# Hilton Hotel Extended Stay Tracker

A lightweight internal web application that helps hotel managers reliably track extended-stay guests and receive **weekly billing reminders** without violating hospitality policies.

## Why This Exists

At Hilton properties, **hospitality and customer-experience policies restrict the use of repeated billing reminders** for extended-stay guests.

That creates a real operational problem:

- Extended-stay guests must be billed on a **weekly cadence**
- Managers are juggling dozens of guests across different check-in dates
- Missing a billing window = revenue leakage
- Manual reminders (notes, spreadsheets, memory) are unreliable

**This app solves that by shifting reminders away from guests and onto the manager.**

The system:

- Tracks extended-stay guests internally
- Calculates how many weeks each guest has stayed
- Sends a **single, guaranteed, once-per-week email reminder** to the manager
- Never contacts the guest directly

No guest spam. No policy violations. No missed billing.

---

## High-Level Architecture

**Backend**

- Go (Gin framework)
- SQLite database
- Background agent running every 24 hours
- Email delivery via Brevo (HTTP API)

**Frontend**

- Plain HTML + JavaScript
- Served directly by the Go backend
- No framework, no build step, no nonsense

**Deployment**

- DigitalOcean Droplet (Linux)
- Go binary managed by `systemd`
- Environment variables stored outside GitHub
- GitHub Actions used for auto-deploy on `main`

---

## How It Works (End-to-End)

1. Manager adds an extended-stay guest through the web UI
2. Guest data is stored in SQLite (`guests` table)
3. A background agent:
    - Runs on startup
    - Runs every 24 hours after that
4. For each guest, the agent:
    - Computes weeks stayed from check-in date
    - Checks if a reminder for that week already exists
    - Sends **exactly one email per guest per week**
    - Records the send in a `notifications` table to prevent duplicates
5. Manager receives a clean billing reminder email with:
    - Guest name
    - Room number
    - Weeks stayed
    - Daily rate
    - Contact info

This design guarantees **idempotency**:

If the server restarts, crashes, or redeploys, reminders are **never duplicated**.

---

## Local Development Setup

### Backend (Go API + Agent)

```bash
go run .

```

- Runs the API on `http://localhost:8080`
- Automatically initializes SQLite (`data.db`)
- Starts the background billing agent

### Frontend (Local Only)

```bash
python3 -m http.server 5500

```

Then open:

```
http://localhost:5500/

```

> ⚠️ When running frontend separately, CORS must be enabled in main.go
> 
> 
> (commented out by default, only for local testing)
> 

---

## Production Deployment (DigitalOcean)

### Environment Variables

Stored **on the server**, not in GitHub:

```bash
cd /opt
nano /etc/extended-stay.env

```

Example:

```
BREVO_API_KEY=your_brevo_api_key
EMAIL_FROM=billing.notifications@yourdomain.com
MANAGER_EMAIL=manager@example.com

```

Reload and restart:

```bash
systemctl daemon-reload
systemctl restart extended-stay

```

Verify:

```bash
systemctl show extended-stay -p Environment |tr' ''\n'

```

---

## GitHub Actions Auto-Deploy

Every push to `main` triggers deployment.

**What happens automatically:**

1. GitHub Action SSHs into the DigitalOcean droplet
2. Pulls the latest code
3. Hard-resets to `origin/main`
4. Restarts the systemd service

```yaml
name:Deploy

on:
push:
branches: ["main"]

jobs:
deploy:
runs-on:ubuntu-latest
steps:
-name:SSHpull+restart
uses:appleboy/ssh-action@v1.0.3
with:
host:${{secrets.DO_HOST}}
username:${{secrets.DO_USER}}
key:${{secrets.DO_SSH_KEY}}
script: |
            set -e
            cd /opt/extended-stay
            git fetch --all
            git reset --hard origin/main
            systemctl restart extended-stay
            systemctl status extended-stay --no-pager -l | head -n 60

```

No manual SSH. No manual restarts. No drift.

---

## Database Design

**guests**

- Stores extended-stay guest details

**notifications**

- Records `(guest_id, period_number)`
- Enforced UNIQUE constraint
- Guarantees one reminder per guest per billing period

This is what makes the system **restart-safe and duplicate-proof**.

---

## Key Design Decisions (Why This Is Solid)

- **SQLite**
    
    Single-tenant, low-traffic internal tool → zero reason for Postgres overhead.
    
- **systemd**
    
    Ensures the service:
    
    - Starts on boot
    - Restarts on failure
    - Runs unattended
- **Background Agent**
    
    No cron jobs, no external schedulers, no race conditions.
    
- **Email to Manager Only**
    
    Complies with hospitality policies while preserving billing accuracy.
    

---

## Who This Is For

- Hotel property managers
- Front desk supervisors
- Operations staff handling extended-stay billing

This is **not** a consumer-facing app.

It’s an internal reliability tool.

---

## Status

✔ Deployed

✔ In production

✔ Restart-safe

✔ Duplicate-proof

✔ Policy-compliant
