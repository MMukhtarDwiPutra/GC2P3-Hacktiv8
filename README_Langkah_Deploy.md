Cara deploy
1. Login terlebih dahulu ke gcloud cli
gcloud auth login

2. Buat project gcloud
gcloud projects create hactiv8-gc2-mmukhtar --name="Hacktiv8-GC2"

3. Set ke project gcloud yang sudah dibuat
gcloud config set project hactiv8-gc2-mmukhtar

4. Build semua docker mejadi image di local
docker build -t gcr.io/hactiv8-gc2-mmukhtar/user-service ./src/user-service
docker build -t gcr.io/hactiv8-gc2-mmukhtar/book-service ./src/book-service
docker build -t gcr.io/hactiv8-gc2-mmukhtar/borrow-service ./src/borrow-service
docker build -t gcr.io/hactiv8-gc2-mmukhtar/api-gateway ./src/api-gateway

5. Upload docker ke gcloud
docker push gcr.io/hactiv8-gc2-mmukhtar/user-service
docker push gcr.io/hactiv8-gc2-mmukhtar/book-service
docker push gcr.io/hactiv8-gc2-mmukhtar/borrow-service
docker push gcr.io/hactiv8-gc2-mmukhtar/api-gateway

6. Run service sesuai dengan port masing masing
gcloud run deploy user-service \
  --image gcr.io/hactiv8-gc2-mmukhtar/user-service \
  --platform managed \
  --region asia-southeast1 \
  --port 50051

gcloud run deploy book-service \
  --image gcr.io/hactiv8-gc2-mmukhtar/book-service \
  --platform managed \
  --region asia-southeast1 \
  --port 50052

gcloud run deploy borrow-service \
  --image gcr.io/hactiv8-gc2-mmukhtar/borrow-service \
  --platform managed \
  --region asia-southeast1 \
  --port 50053

gcloud run deploy api-gateway \
  --image gcr.io/hactiv8-gc2-mmukhtar/api-gateway \
  --platform managed \
  --region asia-southeast1 \
  --port 8080