# Rollback Notes

## before-v6-model-sync-20260620-000142

This snapshot was created before the v6 model-sync, model marketplace, endpoint, and user navigation update.

Local source archive:

```bash
D:/Users/Administrator.LZ-202601111636/Documents/One API 2/backups/sub2api-source-before-v6-model-sync-20260619-233927.zip
```

Server backup directory:

```bash
/home/ubuntu/tokenapifuel-backups/before-v6-model-sync-20260620-000142
```

Server Docker image tag:

```bash
sub2api-oneapi-th:before-v6-model-sync-20260620-000142
```

Server database backup:

```bash
/home/ubuntu/tokenapifuel-backups/before-v6-model-sync-20260620-000142/database-20260620-000142.sql
```

Server data directory backup:

```bash
/home/ubuntu/tokenapifuel-backups/before-v6-model-sync-20260620-000142/data-dir-20260620-000142.tgz
```

Rollback commands:

```bash
docker stop sub2api || true
docker rm sub2api || true

sudo tar -xzf /home/ubuntu/tokenapifuel-backups/before-v6-model-sync-20260620-000142/data-dir-20260620-000142.tgz -C /

docker run -d \
  --name sub2api \
  --restart unless-stopped \
  --network sub2api-deploy_sub2api-network \
  -p 8082:8080 \
  -v /opt/sub2api-deploy/data:/app/data \
  sub2api-oneapi-th:before-v6-model-sync-20260620-000142

cat /home/ubuntu/tokenapifuel-backups/before-v6-model-sync-20260620-000142/database-20260620-000142.sql \
  | docker exec -i sub2api-postgres psql -U sub2api -d sub2api
```

## before-admin-playground-models-20260619-153725

This snapshot was created before the admin/playground/models marketplace upgrade.

Local source archive:

```bash
D:/Users/Administrator.LZ-202601111636/Documents/One API 2/backups/before-admin-playground-models-20260619-153725/sub2api-source-20260619-153725.tgz
```

Server backup directory:

```bash
/home/ubuntu/tokenapifuel-backups/before-admin-playground-models-20260619-153725
```

Server Docker image tag:

```bash
sub2api-oneapi-th:before-admin-playground-models-20260619-153725
```

Server database backup:

```bash
/home/ubuntu/tokenapifuel-backups/before-admin-playground-models-20260619-153725/database-20260619-153725.sql
```

Server data directory backup:

```bash
/home/ubuntu/tokenapifuel-backups/before-admin-playground-models-20260619-153725/data-dir-20260619-153725.tgz
```

Example server rollback commands:

```bash
cd /opt/sub2api-deploy
docker stop sub2api || true
docker rm sub2api || true

sudo tar -xzf /home/ubuntu/tokenapifuel-backups/before-admin-playground-models-20260619-153725/data-dir-20260619-153725.tgz -C /

docker run -d \
  --name sub2api \
  --restart unless-stopped \
  --network sub2api-deploy_sub2api-network \
  -p 8082:8080 \
  -v /opt/sub2api-deploy/data:/app/data \
  sub2api-oneapi-th:before-admin-playground-models-20260619-153725
```

Database restore depends on the active database engine. For the current containerized deployment, prefer restoring the `/opt/sub2api-deploy/data` backup first. Use `database-20260619-153725.sql` if the SQL database was externalized.
