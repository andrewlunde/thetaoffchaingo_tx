```
git add .
git commit -m "query: changes for v0.2.1"
git tag v0.2.1
git push origin v0.2.1
git push
GOPROXY=proxy.golang.org go list -m github.com/andrewlunde/thetaoffchaingo_query@v0.2.1

go get github.com/andrewlunde/thetaoffchaingo_query@v0.2.1

// in go.mod of application using this module
require github.com/thetatoken/query v0.0.0
replace github.com/thetatoken/query v0.0.0 => github.com/andrewlunde/thetaoffchaingo_query v0.2.1

```

