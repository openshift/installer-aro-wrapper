#!/bin/bash
# Usage:
#   export CLUSTER=myloadtest RESOURCEGROUP=myloadtest VERSION=4.21.22
#   bash test.sh create   # provision clusters
#   bash test.sh delete   # delete clusters (VERSION not required)
#
# Region selection prompt:
#   Single region  : enter a number          e.g. 11
#   Multiple regions: enter space-separated  e.g. 1 3 7
#   All regions    : enter 'a'
#
# Each selected region gets a tmux window with $CONCURRENCY panes (default 5).
# Switch windows with Ctrl+b n. Session is named 'loadtest'.

usage() {
	echo -e "usage: ${0} <create|delete>"
	exit 1
}

CLUSTER=${CLUSTER:-loadtest}
RESOURCEGROUP=${RESOURCEGROUP:-loadtest}
CONCURRENCY=${CONCURRENCY:-5}
BASEDIR=$(dirname "$0")

case "$1" in
  create)
    COMMAND="$BASEDIR/create.sh"
    if [ -z $VERSION ]; then
      echo "env VERSION required"
      exit 1
    fi
    ;;
  delete)
    COMMAND="$BASEDIR/delete.sh"
    ;;
  *)
    usage
    exit 1
    ;;
esac

regions=()
while IFS= read -r line; do
  line="${line%$'\r'}"  # strip Windows carriage returns
  [[ -n "$line" ]] && regions+=("$line")
done < "$BASEDIR/regions.txt"
echo "Available regions:"
for i in "${!regions[@]}"; do
  printf "  %3d) %s\n" "$((i+1))" "${regions[$i]}"
done
echo ""
echo "Enter region number(s) space-separated, or 'a' for all regions:"
read -r selection

if [[ "$selection" == "a" || "$selection" == "all" ]]; then
  selected_regions=("${regions[@]}")
else
  selected_regions=()
  for num in $selection; do
    idx=$((num - 1))
    if [[ $idx -ge 0 && $idx -lt ${#regions[@]} ]]; then
      selected_regions+=("${regions[$idx]}")
    else
      echo "Invalid region number: $num"
      exit 1
    fi
  done
fi

if [[ ${#selected_regions[@]} -eq 0 ]]; then
  echo "No regions selected"
  exit 1
fi

# one tmux window per region, each with $CONCURRENCY panes running in parallel
SESSION="loadtest"
tmux start-server
tmux new-session -d -s "$SESSION" -n "${selected_regions[0]}"

for LOCATION in "${selected_regions[@]}"; do
  if [[ "$LOCATION" != "${selected_regions[0]}" ]]; then
    tmux new-window -t "$SESSION" -n "$LOCATION"
  fi

  for (( i=1; i<CONCURRENCY; i++ )); do
    tmux split-window -h -t "$SESSION:$LOCATION"
  done
  tmux select-layout -t "$SESSION:$LOCATION" even-horizontal

  for (( i=0; i<CONCURRENCY; i++ )); do
    tmux send-keys -t "$SESSION:$LOCATION.$i" \
      "LOCATION=$LOCATION CLUSTER=$CLUSTER-$LOCATION-$i RESOURCEGROUP=$RESOURCEGROUP-$LOCATION-$i VERSION=$VERSION $COMMAND" Enter
  done
done

tmux attach-session -t "$SESSION"
