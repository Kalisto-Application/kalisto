import React from 'react';

export const Environments: React.FC = () => {
  return (
    <div className="p-4">
      <p>There are your variables.</p>
      <br></br>
      <p>Currently Kalisto supports only global variables</p>
      <br></br>
      <p>Define a regular JS object</p>
      <p>
        For instance: <code>&#123;&#125;</code>
      </p>
      <br></br>
      <p>It becomes available in your requests as `g` object</p>
    </div>
  );
};
