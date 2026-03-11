import React from 'react';
import { Editor } from '@monaco-editor/react';
import { Wand2 } from 'lucide-react';

interface EditorContainerProps {
  code: string;
  onCodeChange: (value: string) => void;
  language: string;
  theme: string;
  displayLang: string;
  languageVersion?: string;
  onFormat: () => void;
  onEditorMount: (editor: any, monaco: any) => void;
}

export const EditorContainer: React.FC<EditorContainerProps> = ({
  code,
  onCodeChange,
  language,
  theme,
  displayLang,
  languageVersion,
  onFormat,
  onEditorMount,
}) => {
  return (
    <div className="flex-[3] flex flex-col bg-ig-card border border-ig-border rounded-xl overflow-hidden min-h-0">
      <div className="px-4 py-3 bg-card-header border-b border-ig-border flex items-center justify-between shrink-0">
        <div className="text-sm font-semibold flex items-center gap-2">
          <div className="w-8 h-8 rounded-full bg-ig-gradient flex items-center justify-center text-xs font-bold text-white">S</div>
          <span>{displayLang}</span>
          {languageVersion && (
            <span className="text-ig-secondary-text text-xs font-normal">({languageVersion})</span>
          )}
        </div>
        <button 
          onClick={onFormat}
          className="flex items-center gap-1.5 px-4 py-1.5 border border-ig-border rounded-lg font-semibold text-sm hover:opacity-80 transition-opacity"
        >
          <Wand2 size={16} />
          <span className="hidden sm:inline">Format</span>
        </button>
      </div>
      <div className="flex-1 w-full bg-editor">
        <Editor
          height="100%"
          language={language === 'golang' ? 'go' : language}
          value={code}
          theme={theme === 'dark' ? 'code-dark' : 'code-light'}
          onChange={(value) => onCodeChange(value || '')}
          onMount={onEditorMount}
          options={{
            automaticLayout: true,
            fontSize: 14,
            padding: { top: 16 },
            minimap: { enabled: false },
            scrollBeyondLastLine: false,
            lineNumbers: 'on',
            renderLineHighlight: 'all',
            fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
            fontWeight: "400",
            lineHeight: 22,
            cursorSmoothCaretAnimation: "on",
            cursorBlinking: "smooth"
          }}
        />
      </div>
    </div>
  );
};
