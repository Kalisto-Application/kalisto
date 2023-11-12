import { useState } from 'react';
import { useStep } from 'usehooks-ts';
import gIcon from '../../assets/icons/g.svg';
import { models } from '../../wailsjs/go/models';

type RequestList = {
  setActiveRequest: (id: string, nameMethod: string) => void;
  requestId: { id: string; fullNameMet: string };
  fullNameMet: string;
  requestFiles: { [key: string]: models.File[] };
};
const RequestList: React.FC<RequestList> = ({
  setActiveRequest,
  requestId,
  requestFiles,
  fullNameMet,
}) => {
  const setKey = () => {
    const keys = Object.keys(requestFiles);

    keys.forEach((key) => {
      return key;
    });
    for (let index = 0; index < keys.length; index++) {
      const element = keys[index];

      if (fullNameMet === element) return element;
    }
  };

  return (
    <>
      {fullNameMet === setKey() &&
        requestFiles[setKey() || ''].map((it, indx) => {
          return (
            <div
              className="flex items-center gap-2 "
              onClick={() => setActiveRequest(it.name, requestId.fullNameMet)}
              key={indx}
              id={fullNameMet}
            >
              <img src={gIcon} />
              <span className={`${it.id === requestId.id ? 'text-red' : ''}`}>
                {it.name}
              </span>
            </div>
          );
        })}
    </>
  );
};
export default RequestList;
