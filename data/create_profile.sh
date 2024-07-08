# Real avatar base64
#curl -X POST http://localhost:8080/profile -H "Content-Type: application/json" -d '{
#  "name": "Master Vegas",
#  "avatar": "'$(cat avatar_base64.txt)'",
#  "posts": [],
#  "liked_posts": [],
#  "notifications": []
#}'

# Emtpty avatar
curl -X POST http://localhost:8080/profile -H "Content-Type: application/json" -d '{
  "name": "Master Vegas",
  "avatar": "",
  "posts": [],
  "liked_posts": [],
  "notifications": []
}'