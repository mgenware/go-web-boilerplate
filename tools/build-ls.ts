import { promises as fsPromises } from 'fs';
import * as camelCase from 'camelcase';

const typeName = 'Dictionary';

function pascalCase(s: string): string {
  return camelCase(s, { pascalCase: true });
}

(async () => {
  let result =
    '/** DO NOT EDIT. This file was automatically generated by `/tools/build-ls.ts`. */\n\npackage localization\n\n';
  result += `// ${typeName} contains entries defined in app localized strings.\n`;
  result += `type ${typeName} struct {\n`;

  // Read the contents of a localized JSON.
  const json = await fsPromises.readFile(
    '../localization/langs/en.json',
    'utf8',
  );
  const ls = JSON.parse(json);
  Object.keys(ls)
    .sort()
    .forEach((key) => {
      result += `\t${pascalCase(key)} string \`json:"${key}"\`\n`;
    });

  result += '}\n';

  await fsPromises.writeFile(
    '../src/app/handler/localization/dictionary.go',
    result,
  );
})();
