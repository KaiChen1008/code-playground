import React from 'react';

interface ConsoleOutputProps {
  output: string;
  snippetId: string | null;
}

export const ConsoleOutput: React.FC<ConsoleOutputProps> = ({ output, snippetId }) => {
  return (
    <div className="flex-1 flex flex-col bg-ig-card border border-ig-border rounded-xl overflow-hidden min-h-[120px] sm:min-h-[150px]">
      <div className="px-4 py-3 bg-card-header border-b border-ig-border flex items-center justify-between shrink-0">
        <div className="text-sm font-semibold">Console Output</div>
        {snippetId && (
          <div className="text-sm text-ig-secondary-text">
            Snippet ID: <span className="text-ig-blue font-semibold">{snippetId}</span>
          </div>
        )}
      </div>
      <pre className="p-4 m-0 font-mono text-xs sm:text-sm text-output whitespace-pre-wrap overflow-y-auto bg-editor flex-1">
        {output}
      </pre>
    </div>
  );
};
