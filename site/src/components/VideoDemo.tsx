export default function VideoDemo() {
  return (
    <section className="px-6 py-24 border-t border-zinc-800/50">
      <div className="max-w-4xl mx-auto">
        <h2 className="text-2xl sm:text-3xl font-bold text-white text-center mb-4">
          See it in action
        </h2>
        <p className="text-zinc-500 text-center mb-12 max-w-xl mx-auto">
          From bare server to live WordPress site in one command.
        </p>

        <div className="rounded-xl border border-zinc-800 overflow-hidden shadow-2xl shadow-black/50">
          <div className="relative w-full" style={{ paddingBottom: "56.25%" }}>
            <iframe
              className="absolute inset-0 w-full h-full"
              src="https://www.youtube.com/embed/oZmqxB_eMMY"
              title="pressctl demo"
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              allowFullScreen
            />
          </div>
        </div>
      </div>
    </section>
  );
}
