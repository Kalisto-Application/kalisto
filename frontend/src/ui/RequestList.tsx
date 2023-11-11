import gIcon from '../../assets/icons/g.svg';
import { dataRequestList } from '../components/MethodCollectionView';

type RequestList = {
  files: dataRequestList[];
  setActiveRequest: (id: string) => void;
  requestId: string;
};
const RequestList: React.FC<RequestList> = ({
  files,
  setActiveRequest,
  requestId,
}) => {
  return (
    <>
      {files.map((it, indx) => (
        <div
          className="flex items-center gap-2 "
          onClick={() => setActiveRequest(it.id)}
          key={indx}
        >
          <img src={gIcon} />
          <span className={`${it.id === requestId ? 'text-red' : ''}`}>
            {it.name}
          </span>
        </div>
      ))}
    </>
  );
};
export default RequestList;
