import classNames from "classnames";
import { useRouter } from "next/router";

function AppLayout({ children, bgConf }) {
  const router = useRouter()

  return (
    <div
      className={classNames({
        "min-h-screen": true,
        "bg-cover": true,
        "bg-no-repeat": true,
        "bg-center": true,
        "bg-gray-900": true,
        [bgConf.bg]: router.query.id === undefined,
      })}
      id="main"
    >
      {children}
    </div>
  );
}

export default AppLayout;
