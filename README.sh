#!/bin/bash

README_FILE="README.md"

temp_file=$(mktemp)

# State tracking
in_section=false
section_type=""

while IFS= read -r line; do
  # When not in a section, look for a start marker
  if [[ "$in_section" = false ]]; then
    if [[ $line =~ ^\<\!--BENCHMARK:([^>]+)--\> ]]; then
      # We found a benchmark start marker
      benchmark_file="${BASH_REMATCH[1]}"
      in_section=true
      section_type="benchmark"

      # Write the marker
      echo "$line" >> "$temp_file"

      # Add benchmark content
      if [[ -f "$benchmark_file" ]]; then
        echo '```' >> "$temp_file"
        cat "$benchmark_file" >> "$temp_file"
        echo '```' >> "$temp_file"
      else
        echo "Warning: Benchmark file $benchmark_file not found" >&2
      fi

    elif [[ $line =~ ^\<\!--COMMAND:([^>]+)--\> ]]; then
      # We found a command start marker
      command="${BASH_REMATCH[1]}"
      in_section=true
      section_type="command"

      # Write the marker
      echo "$line" >> "$temp_file"

      # Add command content
      echo '```' >> "$temp_file"
      echo "$ $command" >> "$temp_file"
      eval "$command" >> "$temp_file"
      echo '```' >> "$temp_file"

    elif [[ $line =~ ^\<\!--SQL:([^>]+)--\> ]]; then
      # We found a sql start marker
      sql_file="${BASH_REMATCH[1]}"
      in_section=true
      section_type="sql"

      # Write the marker
      echo "$line" >> "$temp_file"

      # Add sql content
      if [[ -f "$sql_file" ]]; then
        echo '```sql' >> "$temp_file"
        cat "$sql_file" >> "$temp_file"
        echo '```' >> "$temp_file"
      else
        echo "Warning: SQL file $sql_file not found" >&2
      fi

    else
      # Regular line, just copy it
      echo "$line" >> "$temp_file"
    fi

  # When in a section, look for the corresponding end marker
  else
    if [[ "$section_type" = "benchmark" && $line =~ ^\<\!--END_BENCHMARK--\> ]]; then
      # Found the end benchmark marker
      echo "$line" >> "$temp_file"
      in_section=false
      section_type=""

    elif [[ "$section_type" = "command" && $line =~ ^\<\!--END_COMMAND--\> ]]; then
      # Found the end command marker
      echo "$line" >> "$temp_file"
      in_section=false
      section_type=""

    elif [[ "$section_type" = "sql" && $line =~ ^\<\!--END_SQL--\> ]]; then
      # Found the end sql marker
      echo "$line" >> "$temp_file"
      in_section=false
      section_type=""

    else
      # Skip this line as we're between markers
      continue
    fi
  fi
done < "$README_FILE"

# Replace the original README with the new one
mv "$temp_file" "$README_FILE"
