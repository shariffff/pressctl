import Hero from "@/components/Hero";
import CommandShowcase from "@/components/CommandShowcase";
import Features from "@/components/Features";
import HowItWorks from "@/components/HowItWorks";
import Footer from "@/components/Footer";

export default function Home() {
  return (
    <main className="min-h-screen">
      <Hero />
      <CommandShowcase />
      <Features />
      <HowItWorks />
      <Footer />
    </main>
  );
}
