go test -run=NOTEST -benchtime=10000x -bench=. -count=20 > x.txt
benchstat x.txt