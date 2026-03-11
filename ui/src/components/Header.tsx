import React from 'react';
import { Play, Trash2, Moon, Sun } from 'lucide-react';
import { Language } from '../api';

interface HeaderProps {
  theme: string;
  onToggleTheme: () => void;
  languages: Language[];
  currentLang: string;
  onLanguageChange: (lang: string) => void;
  onRun: () => void;
  onDelete: () => void;
  snippetId: string | null;
  isRunning: boolean;
}

export const Header: React.FC<HeaderProps> = ({
  theme,
  onToggleTheme,
  languages,
  currentLang,
  onLanguageChange,
  onRun,
  onDelete,
  snippetId,
  isRunning,
}) => {
  return (
    <header className="w-full px-5 py-3 flex justify-between items-center border-b border-ig-border bg-header z-50">
      <div 
        className="text-2xl font-bold bg-ig-gradient bg-clip-text text-transparent cursor-pointer tracking-tighter"
        onClick={() => window.location.href = '/'}
      >
        code
      </div>
      <div className="flex gap-4 items-center">
        <select 
          id="language"
          className="bg-ig-bg text-ig-text border border-ig-border px-3 py-1.5 rounded-lg text-sm outline-none cursor-pointer"
          value={currentLang}
          onChange={(e) => onLanguageChange(e.target.value)}
        >
          {languages.map(lang => (
            <option key={lang.name} value={lang.name}>
              {lang.name.charAt(0).toUpperCase() + lang.name.slice(1)}
            </option>
          ))}
        </select>
        <button 
          onClick={onToggleTheme}
          className="flex items-center gap-1.5 px-4 py-1.5 border border-ig-border rounded-lg font-semibold text-sm hover:opacity-80 transition-opacity"
          title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
        >
          {theme === 'dark' ? <Sun size={16} /> : <Moon size={16} />}
          <span className="hidden sm:inline">{theme === 'dark' ? 'Light' : 'Dark'}</span>
        </button>
        <button 
          onClick={onRun}
          disabled={isRunning}
          className="flex items-center gap-1.5 px-4 py-1.5 bg-ig-blue text-white rounded-lg font-semibold text-sm hover:opacity-80 transition-opacity disabled:opacity-50"
        >
          <Play size={16} fill="currentColor" />
          <span className="hidden sm:inline">Run</span>
        </button>
        {snippetId && (
          <button 
            onClick={onDelete}
            className="flex items-center gap-1.5 px-4 py-1.5 border border-ig-border text-ig-red rounded-lg font-semibold text-sm hover:opacity-80 transition-opacity"
          >
            <Trash2 size={16} />
          </button>
        )}
      </div>
    </header>
  );
};
