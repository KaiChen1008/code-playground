import { useState, useEffect, useRef } from 'react';
import * as api from './api';
import { samples, extensions } from './constants';
import { Header } from './components/Header';
import { EditorContainer } from './components/EditorContainer';
import { ConsoleOutput } from './components/ConsoleOutput';

function App() {
  const [theme, setTheme] = useState(() => localStorage.getItem('theme') || 'light');
  const [languages, setLanguages] = useState<api.Language[]>([]);
  const [currentLang, setCurrentLang] = useState('golang');
  const [code, setCode] = useState(samples['golang']);
  const [output, setOutput] = useState('Welcome to code. Write code and press Run.');
  const [snippetId, setSnippetId] = useState<string | null>(null);
  const [originalCode, setOriginalCode] = useState(samples['golang']);
  const [isRunning, setIsRunning] = useState(false);
  const editorRef = useRef<any>(null);

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  }, [theme]);

  useEffect(() => {
    api.fetchLanguages().then(setLanguages).catch(console.error);

    const pathId = window.location.pathname.substring(1);
    if (pathId && pathId.length === 6) {
      api.loadSnippet(pathId).then(data => {
        setSnippetId(data.id);
        setCurrentLang(data.language);
        setCode(data.code);
        setOriginalCode(data.code);
        if (data.output) setOutput(data.output);
      }).catch(console.error);
    }
  }, []);

  const toggleTheme = () => setTheme(prev => prev === 'dark' ? 'light' : 'dark');

  const handleLanguageChange = (lang: string) => {
    setCurrentLang(lang);
    setSnippetId(null);
    const newCode = samples[lang] || "";
    setCode(newCode);
    setOriginalCode(newCode);
    window.history.pushState(null, '', '/');
    setOutput('Welcome to code. Write code and press Run.');
  };

  const handleRun = async () => {
    setIsRunning(true);
    setOutput('Running...');
    try {
      const data = await api.runCode(currentLang, code, snippetId && code === originalCode ? snippetId : undefined);
      if (data.error) {
        setOutput('Error: ' + data.error);
      } else {
        setOutput(data.output);
        if (data.id) {
          setSnippetId(data.id);
          setOriginalCode(code);
          window.history.pushState(null, '', `/${data.id}`);
        }
      }
    } catch (err: any) {
      setOutput('Error: ' + err.message);
    } finally {
      setIsRunning(false);
    }
  };

  const handleFormat = async () => {
    if (currentLang === 'golang') {
      try {
        const data = await api.formatCode(currentLang, code);
        if (data.code) {
          setCode(data.code);
        }
      } catch (err) {
        console.error('Format error:', err);
      }
    } else if (editorRef.current) {
      editorRef.current.getAction('editor.action.formatDocument').run();
    }
  };

  const handleDelete = async () => {
    if (!snippetId) return;
    if (!confirm('Are you sure you want to delete this snippet?')) return;
    try {
      await api.deleteSnippet(snippetId);
      window.location.href = '/';
    } catch (err: any) {
      alert('Error: ' + err.message);
    }
  };

  const handleEditorDidMount = (editor: any, monaco: any) => {
    editorRef.current = editor;

    monaco.editor.defineTheme('code-light', {
      base: 'vs',
      inherit: true,
      rules: [
        { token: 'keyword', foreground: 'd73a49', fontStyle: 'bold' },
        { token: 'comment', foreground: '6a737d', fontStyle: 'italic' },
        { token: 'string', foreground: '032f62' },
        { token: 'number', foreground: '005cc5' },
        { token: 'type', foreground: '6f42c1' },
        { token: 'function', foreground: '005cc5' }
      ],
      colors: {
        'editor.background': '#ffffff',
        'editor.lineHighlightBackground': '#f6f8fa',
        'editorCursor.foreground': '#24292e',
        'editorIndentGuide.background': '#e1e4e8',
        'editor.selectionBackground': '#0366d625'
      }
    });

    monaco.editor.defineTheme('code-dark', {
      base: 'vs-dark',
      inherit: true,
      rules: [
        { token: 'keyword', foreground: 'ff7b72', fontStyle: 'bold' },
        { token: 'comment', foreground: '8b949e', fontStyle: 'italic' },
        { token: 'string', foreground: 'a5d6ff' },
        { token: 'number', foreground: '79c0ff' },
        { token: 'type', foreground: 'd2a8ff' },
        { token: 'function', foreground: '79c0ff' }
      ],
      colors: {
        'editor.background': '#131722',
        'editor.lineHighlightBackground': '#1c2230',
        'editorCursor.foreground': '#f0f6fc',
        'editorIndentGuide.background': '#2d3444',
        'editor.selectionBackground': '#388bfd40'
      }
    });

    monaco.editor.setTheme(theme === 'dark' ? 'code-dark' : 'code-light');

    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, handleFormat);
  };

  const languageVersion = languages.find(l => l.name === currentLang)?.version;
  const displayLang = `main.${extensions[currentLang] || 'txt'}`;

  return (
    <div className="flex flex-col h-full w-full bg-ig-bg text-ig-text">
      <Header 
        theme={theme}
        onToggleTheme={toggleTheme}
        languages={languages}
        currentLang={currentLang}
        onLanguageChange={handleLanguageChange}
        onRun={handleRun}
        onDelete={handleDelete}
        snippetId={snippetId}
        isRunning={isRunning}
      />

      <main className="flex-1 w-full flex flex-col p-3 sm:p-5 gap-3 sm:gap-5 overflow-hidden min-h-0">
        <EditorContainer 
          code={code}
          onCodeChange={setCode}
          language={currentLang}
          theme={theme}
          displayLang={displayLang}
          languageVersion={languageVersion}
          onFormat={handleFormat}
          onEditorMount={handleEditorDidMount}
        />

        <ConsoleOutput 
          output={output}
          snippetId={snippetId}
        />
      </main>
    </div>
  );
}

export default App;
