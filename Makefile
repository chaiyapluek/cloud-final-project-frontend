templ:
	templ generate

tailwindcss:
	npx tailwindcss --config tailwind.config.js -i input.css -o static/css/style.css

generate: tailwindcss templ

run: generate	
	go run ./src/cmd