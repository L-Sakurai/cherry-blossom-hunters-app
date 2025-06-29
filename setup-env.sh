# railway login
# railway link

while IFS='=' read -r key value; do
  if [[ -n "$key" && ! "$key" =~ ^# ]]; then
    echo "Setting $key..."
    railway environment set "$key" "$value"
  fi
done < .env_source

