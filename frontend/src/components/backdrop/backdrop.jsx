function Backdrop({ onClick }) {
    return <div onClick={onClick} className="fixed bg-black opacity-30 z-40 top-0 left-0 w-full h-full" />;
}

export default Backdrop