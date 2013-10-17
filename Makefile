extension:
			rm -rf build/extension
			cp -r extension build/extension
			coffee --compile build/extension/*.coffee
			rm -rf build/extension/*.coffee
			zip -r build/extension.zip build/extension/*

.PHONY: extension
