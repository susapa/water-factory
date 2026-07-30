# Deploying to a VPS (Docker Compose)

Plain HTTP deploy — no domain/SSL required. Point a browser at `http://<vps-ip>/`.

## 0. Provision the VPS

Pick a provider and spin up an instance — for dev-only use, the cheapest tier is fine:

- **Vultr / DigitalOcean** — ~$4-6/mo, 1 vCPU / 1-2GB RAM, billed hourly (destroy when not needed)
- **Hetzner CX22** — ~€4/mo, 2 vCPU / 4GB RAM, best specs-per-baht but EU-only datacenters
- **Oracle Cloud Free Tier** — free forever (ARM, 4 vCPU / 24GB RAM) if you can get the signup through

When creating the instance:
- OS: **Ubuntu 24.04 LTS** (or 22.04)
- Auth: **SSH key** (paste your public key, e.g. `~/.ssh/id_ed25519.pub` — generate one with `ssh-keygen -t ed25519` if you don't have one) instead of a password
- Note the VPS's public IP once it's created

### First login + basic hardening

```bash
ssh root@<vps-ip>
```

Create a non-root sudo user (don't run Docker/the app as root long-term):

```bash
adduser deploy
usermod -aG sudo deploy
rsync --archive --chown=deploy:deploy ~/.ssh /home/deploy
```

Log out, then back in as the new user:

```bash
ssh deploy@<vps-ip>
```

Lock down SSH — edit `/etc/ssh/sshd_config`, set:

```
PasswordAuthentication no
PermitRootLogin no
```

Then `sudo systemctl restart ssh`.

Enable the firewall (SSH + HTTP only — see the Firewall section below for the full command):

```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw enable
```

## 1. Install Docker on the VPS (Ubuntu)

```bash
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

Log out/in (or `newgrp docker`) so the docker group membership takes effect.

## 2. Get the code

```bash
git clone https://github.com/susapa/water-factory.git
cd water-factory
```

## 3. Configure environment

```bash
cp .env.production.example .env
nano .env   # fill in DB_PASSWORD and JWT_SECRET with real random values
```

Generate strong values with:

```bash
openssl rand -base64 48
```

**Never reuse the dev secrets from `backend/.env`.**

## 4. Deploy

```bash
chmod +x deploy.sh
./deploy.sh
```

This builds the backend (Go) and frontend (Angular + nginx) images and starts all three services (postgres, backend, frontend). Database migrations run automatically on backend startup.

## 5. Verify

- Frontend: `http://<vps-ip>/`
- Login with the seeded admin: `admin@water.local` / `Admin@1234` (change this password after first login)
- Backend health check (only reachable inside the VPS, since port 8080 isn't published — see Firewall below): `curl http://localhost:8080/health` on the VPS itself, or `docker compose logs backend`

## Operations

```bash
docker compose logs -f backend     # tail backend logs
docker compose logs -f frontend    # tail frontend/nginx logs
docker compose ps                  # service status
docker compose restart backend     # restart one service
docker compose down                # stop everything (data persists in the postgres_data volume)
```

## Updating after a code change

```bash
git pull
./deploy.sh
```

## Optional: load demo/test seed data

```bash
docker compose exec -T postgres psql -U water_user -d water_factory < backend/seed_demo_data.sql
```

## Firewall

Only these ports need to be open to the internet:

- `22` — SSH
- `80` — frontend (nginx serves the SPA and proxies `/api/` to the backend)

`8080` (backend) and `5432` (postgres) are **not published** in `docker-compose.yml` — they're only reachable between containers over the internal `water_net` network. Note: Docker manages its own iptables rules for any port it *does* publish, which bypass `ufw` — so don't add `ports:` mappings back for these without also restricting them (e.g. bind to `127.0.0.1:8080:8080` instead of `8080:8080`, or use `ufw-docker`).

```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw enable
```

## Security notes before deploying for real

- `backend/.env` and root `.env` are git-ignored and not committed — good, keep it that way.
- Generate a fresh `JWT_SECRET` and `DB_PASSWORD` for the VPS; don't copy your local dev `.env`.
- Change the seeded admin password (`admin@water.local`) immediately after first login.
