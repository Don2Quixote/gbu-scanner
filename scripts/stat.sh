# Pride for bash knowledges.

# Counts all project's packages..
PACKAGES_COUNT=$(go list ./... | wc -l)

# Counts all .go files in project.
FILES_COUNT=$(find | grep "\\.go$" | wc -l)

# Coutns all lines and chars in all .go files in project.
# Result is space separated: lines chars.
WC_RESULT=$(cat $(find | grep "\\.go$") | wc -l -c)

# I don't know why does it create array from string with spaces.
# (Are there really ones that love bash?)
WC_RESULT_ARR=( $WC_RESULT )

LINES_COUNT=${WC_RESULT_ARR[0]}
CHARS_COUNT=${WC_RESULT_ARR[1]}

# Output collected information.
printf "%8s: %d\n" Packages $PACKAGES_COUNT
printf "%8s: %d\n" Files $FILES_COUNT
printf "%8s: %d\n" Lines $LINES_COUNT
printf "%8s: %d\n" Chars $CHARS_COUNT
