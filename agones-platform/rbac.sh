#!/bin/bash

export CLOUDBUILD_ACCOUNT="827464150398@cloudbuild.gserviceaccount.com"
# echo "Ensure that CloudBuild has access to project"
# gcloud projects add-iam-policy-binding "$PROJECT_ID" \
#     --member serviceAccount:"$CLOUDBUILD_ACCOUNT" \
#     --project "$PROJECT_ID" \
#     --role roles/editor || err

echo "Add Cloudbuild to cluster if not already there"
gcloud beta container clusters get-credentials agones-cluster --region australia-southeast1-a --project eloquent-hold-251023 || true

kubectl create clusterrolebinding admin-gke-binding --clusterrole=cluster-admin --user="sai.maduri@kasna.com.au" || true

#kubectl create clusterrolebinding cloudbuild-gke-binding --clusterrole=cluster-admin --user=$CLOUDBUILD_ACCOUNT || true
kubectl --namespace kube-system create serviceaccount tiller

kubectl create clusterrolebinding tiller-cluster-rule \
 --clusterrole=cluster-admin --serviceaccount=kube-system:tiller

helm init --service-account tiller