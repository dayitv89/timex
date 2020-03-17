# timex

# Run timeout test

\$ `go test ./timeout -count=1 -v -coverprofile /var/tmp/timeout.out`

see the covereage report as
\$ `go tool cover -html=/var/tmp/timeout.out`
