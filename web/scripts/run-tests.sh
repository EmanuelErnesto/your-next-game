#!/bin/bash

# Find all test files and run them serially in isolated processes to prevent shared global DOM pollution
TEST_FILES=$(find src -name "*.test.tsx" -o -name "*.test.ts")
EXIT_CODE=0

echo "🚀 Iniciando suíte de testes unitários do frontend..."
echo ""

for file in $TEST_FILES; do
  echo "--------------------------------------------------"
  echo "🧪 Rodando teste: $file"
  bun test --preload ./tests/setup.ts "$file"
  
  if [ $? -ne 0 ]; then
    echo "❌ Falha no teste: $file"
    EXIT_CODE=1
  else
    echo "✅ Passou: $file"
  fi
done

echo "--------------------------------------------------"
if [ $EXIT_CODE -eq 0 ]; then
  echo "🎉 Todos os testes passaram com sucesso!"
else
  echo "😢 Algum teste falhou. Por favor, verifique os erros acima."
fi

exit $EXIT_CODE
