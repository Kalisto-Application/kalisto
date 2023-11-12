import { useEffect, useState } from 'react';
import { useStep } from 'usehooks-ts';
import gIcon from '../../assets/icons/g.svg';
import { models } from '../../wailsjs/go/models';

type RequestList = {
  setActiveRequest: (id: string, nameMethod: string) => void;
  requestId: { id: string; fullNameMet: string };
  fullNameMet: string;
  requestFiles: { [key: string]: models.File[] };
  setKey: (fullNameMet: string) => string;
};
const RequestList: React.FC<RequestList> = ({
  requestFiles,
  fullNameMet,
  setKey,
}) => {
  return (
    <>
      {requestFiles[setKey(fullNameMet)].map((it, indx) => {
        return (
          <div className="flex items-center gap-2 " key={indx}>
            <img src={gIcon} />
            <span>{it.name}</span>
          </div>
        );
      })}
    </>
  );
};
export default RequestList;
