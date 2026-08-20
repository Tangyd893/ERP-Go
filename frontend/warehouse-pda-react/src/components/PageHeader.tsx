import { ArrowLeft } from "lucide-react";
import { useNavigate } from "react-router-dom";

interface PageHeaderProps {
  title: string;
  backTo?: string;
  onBack?: () => void;
  right?: React.ReactNode;
}

export default function PageHeader({ title, backTo, onBack, right }: PageHeaderProps) {
  const navigate = useNavigate();

  function handleBack() {
    if (onBack) {
      onBack();
      return;
    }
    if (backTo) {
      navigate(backTo);
      return;
    }
    navigate(-1);
  }

  return (
    <div className="mb-3 flex items-center justify-between gap-2">
      <div className="flex min-w-0 flex-1 items-center gap-2">
        <button
          type="button"
          className="pda-touch-target flex shrink-0 items-center justify-center rounded-md text-pda-text-secondary"
          onClick={handleBack}
        >
          <ArrowLeft className="h-5 w-5" />
        </button>
        <h2 className="truncate text-base font-semibold text-pda-text">{title}</h2>
      </div>
      {right ? <div className="flex shrink-0 items-center gap-2">{right}</div> : null}
    </div>
  );
}
