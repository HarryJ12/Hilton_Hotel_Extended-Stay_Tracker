# Hilton Hotel Extended Stay Tracker

A lightweight internal web application that helps hotel managers reliably track extended-stay guests and receive **weekly billing reminders**.

## Usage

At Hilton properties, **hospitality and customer-experience policies restrict the use of repeated billing reminders** for extended-stay guests.

Extended-stay guests **must be billed on time**, but this policy causes a few issues:

- Managers are juggling dozens of guests across different check-in dates
- Missing a billing window leads to delayed or inaccurate billing reporting
- Manual reminders (notes, spreadsheets, memory) are unreliable

**This app solves that by shifting reminders away from guests and onto the manager**:

- Tracks extended-stay guests internally after guest details are manually entered
- Calculates how many weeks each guest has stayed
- Sends a **single, guaranteed, once-per-week email reminder** to the manager

No policy violations and no missed billing.

## Architecture

### Backend

- **Go (Gin framework)** handles all HTTP routing, request validation, and API responses. Also serves the frontend files directly.
    
- **SQLite Database** stores guest records and notification history. A unique constraint ensures billing reminders are sent once per guest per billing period, even across restarts.
    
- **Background Agent (24-hour interval)** runs automatically on startup and every 24 hours thereafter. The agent:
    - Calculates how long each guest has stayed
    - Determines whether a billing reminder is due
    - Sends a reminder only if one has not already been sent for that period
- **Email Delivery (Brevo HTTP API)** sends billing reminders to the manager using Brevo’s transactional email API over HTTP. 

### Frontend

- Minimal **HTML + JavaScript** frontend used by staff to add, view, and remove extended-stay guests.

    - Static frontend files are served directly by the Go application
    
### Deployment

- **DigitalOcean Droplet (Ubuntu Linux)**
    
- **systemd-managed Go Binary** allowing:    
    - Automatic startup on boot
    - Automatic restarts on failure
    - Unattended, long-running operation

- **Environment-Based Configuration** ensures sensitive configuration is loaded from a secure environment file and attached to the `systemd` service so it is not exposed.
    - `/etc/extended-stay.env` (Example environment file):
    
        ```bash
        BREVO_API_KEY=your_key_here
        EMAIL_FROM=your_email_here
        MANAGER_EMAIL=your_email_here
        ```

- **GitHub Actions (Auto-Deploy):** Every push to the `main` branch triggers a deployment workflow that:
    
    - SSHs into the DigitalOcean droplet
    - Pulls the latest code
    - Restarts the `systemd` service

This setup provides a simple, secure, and restart-safe deployment pipeline

## Database Design

**guests**

- Stores extended-stay guest information used for billing and tracking

**notifications**

- Tracks which billing reminders have already been sent
- Uses `(guest_id, period_number)` with a `UNIQUE` constraint to prevent duplicate reminders

## How It Works

1. Manager or staff adds an extended-stay guest through the web UI
2. Guest information is persisted in SQLite
3. A background agent runs:
    - On application startup
    - After a new guest is added
    - Every 24 hours
4. For each guest, the agent:
    - Calculates weeks stayed from the check-in date
    - Determines whether a reminder has already been sent for the current period
    - Sends **exactly one billing reminder per guest per week**
5. The manager receives an email containing:
    - Guest name
    - Room number
    - Weeks stayed
    - Daily rate
    - Contact information

## Local Development Setup

### Environment Variables

Add the environment variables to your `~/.zshrc` file so they are automatically loaded in every terminal session.

```bash
nano ~/.zshrc

```

Add the following lines:

```bash
export BREVO_API_KEY=your_key_here
export EMAIL_FROM=your_email_here
export MANAGER_EMAIL=your_email_here

```

Save and reload:

```bash
source ~/.zshrc

```

Verify:

```bash
echo$BREVO_API_KEY

```

### Backend (Go API + Agent)

From the project root:
```bash
cd /backend
go run .

```

### Frontend

In a new terminal, from the project root:
```bash
cd /frontend
python3 -m http.server 5500

```

Then open in browser:

```bash
http://localhost:5500/
```
> ⚠️ When running frontend locally, certain configuration changes are required
> 
> 
> (commented out by default, only for local testing)
> 
