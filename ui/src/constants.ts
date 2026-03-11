export const samples: Record<string, string> = {
  python: 'print("Hello, code!")',
  javascript: 'console.log("Hello, code!");',
  golang: 'package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello, code!")\n}',
  rust: 'fn main() {\n    println!("Hello, code!");\n}',
  cpp: '#include <iostream>\n\nint main() {\n    std::cout << "Hello, code!" << std::endl;\n    return 0;\n}'
};

export const extensions: Record<string, string> = {
  golang: 'go',
  python: 'py',
  javascript: 'js',
  rust: 'rs',
  cpp: 'cpp'
};
