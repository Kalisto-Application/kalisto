import React, { useState } from 'react';
import { FindProtoFiles, CreateWorkspace } from './../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import Button from '../ui/Button';
import Popup from '../ui/Popup';
import upload from '../../assets/icons/upload.svg';
import folder from '../../assets/icons/folder.svg';
import close from '../../assets/icons/close.svg';
import file from '../../assets/icons/file.svg';

// State для Popup
//   const [isOpenCreateWorkspace, setIsOpenCreateWorkspace] = useState(false);

// кнопка для вызова Popup

// <button onClick={() => setIsOpenDeletePopup(true)}>open delete</button>;

// компонента для Popup
{
  /* <CreateWorkspacePopup
        open={isOpenCreateWorkspace}
        onClose={() => setIsOpenCreateWorkspace(false)}
      /> */
}

interface propsCreateWorkspacePopup {
  onClose: () => void;
  open: boolean;
}
interface propsCreateWorkspace {
  onClose: () => void;
}

interface propsUpload {
  data: models.ProtoDir;
  setIsDisabled: (value: { inpitDisabled: boolean; upload: boolean }) => void;
  setData: (dir: models.ProtoDir) => void;
  clearState: () => void;
  isDisabled: {
    inpitDisabled: boolean;
    upload: boolean;
  };
}

const CreateWorkspacePopup: React.FC<propsCreateWorkspacePopup> = ({
  onClose,
  open,
}) => {
  return (
    <Popup onClose={onClose} isOpen={open} title="Name of New Request">
      <CreateWorkspaceComponets onClose={onClose} />
    </Popup>
  );
};

export default CreateWorkspacePopup;

const CreateWorkspaceComponets: React.FC<propsCreateWorkspace> = ({
  onClose,
}) => {
  const [isDisabled, setIsDisabled] = useState({
    inpitDisabled: true,
    upload: true,
  });

  const [valueInp, setValueInp] = useState('');

  const [data, setData] = useState(
    new models.ProtoDir({
      folder: '',
      files: [],
    })
  );

  const updateTextValueInp = (e: React.FormEvent<HTMLInputElement>) => {
    let newValue = e.currentTarget.value;
    setValueInp(newValue);
    if (newValue.length >= 1) {
      setIsDisabled({
        ...isDisabled,
        inpitDisabled: false,
      });
    }
    if (newValue.length < 1) {
      setIsDisabled({
        ...isDisabled,
        inpitDisabled: true,
      });
    }
  };
  const getPropDisabled = () => {
    if (!isDisabled.inpitDisabled && !isDisabled.upload) {
      return false;
    }
    return true;
  };

  return (
    <div className="grid">
      <div className="mb-10 grid gap-y-2">
        <Upload
          data={data}
          setData={setData}
          clearState={() =>
            setData({
              folder: '',
              files: [],
            })
          }
          setIsDisabled={(obj) => setIsDisabled(obj)}
          isDisabled={isDisabled}
        />
        <label className="leading-6">Name</label>
        <input
          value={valueInp}
          onChange={(e) => updateTextValueInp(e)}
          className="rounded-md border-[1px] border-borderFill bg-primaryFill px-4 py-1.5 placeholder-secondaryText"
          type="text"
          placeholder="Name of New Request"
        />
      </div>
      <div className="justify-self-end">
        <Button
          text="Cancel"
          onClick={onClose}
          className="mr-[10px] border-[1px] border-[#343434]"
        />
        <Button
          disabled={getPropDisabled()}
          text="Create"
          onClick={() => {}}
          className="bg-primaryGeneral text-lg font-medium"
        />
      </div>
    </div>
  );
};

const Upload: React.FC<propsUpload> = ({
  data,
  clearState,
  setData,
  setIsDisabled,
  isDisabled,
}) => {
  const findFiles = () => {
    FindProtoFiles().then((res) => {
      if (res.files.length !== 0 && res.folder !== '') {
        setIsDisabled({
          ...isDisabled,
          upload: false,
        });

        setData({
          folder: res.folder,
          files: res.files,
        });
      }
    });
  };

  return (
    <div>
      {data.files.length === 0 ? (
        <button
          onClick={findFiles}
          className="mb-6 flex min-h-[179px] w-full flex-col items-center justify-center border-2 border-dashed border-borderFill"
        >
          <img src={upload} className="mb-6" alt="upload" />
          <span className="font-sm  font-['Roboto_Mono'] font-bold text-secondaryText">
            Choose folder to upload
          </span>
          <span className="text-xs  text-[#5E5E5E]">proto</span>
        </button>
      ) : (
        <div className="max mb-6 max-h-52 overflow-auto rounded-md bg-textBlockFill px-5 py-3.5 ">
          <div className="flex justify-between [&:not(:last-child)]:mb-1">
            <div className="flex items-center gap-x-2">
              <img src={folder} alt="file" />
              <span className="text-sm text-blueTextPath">{data.folder}</span>
            </div>
            <button
              onClick={() => {
                setIsDisabled({
                  ...isDisabled,
                  upload: true,
                });
                clearState();
              }}
            >
              <img src={close} alt="clearState" />
            </button>
          </div>
          <ul>
            {data.files.map((f, indx) => {
              return (
                <li
                  key={indx}
                  className="flex gap-x-2  pl-5 text-sm [&:not(:last-child)]:mb-1"
                >
                  <img src={file} alt="file" />
                  <span>{f}</span>
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </div>
  );
};
