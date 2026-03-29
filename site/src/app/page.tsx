import Hero from "@/components/Hero";
import CommandShowcase from "@/components/CommandShowcase";
import Features from "@/components/Features";
import Architecture from "@/components/Architecture";
import HowItWorks from "@/components/HowItWorks";
import Footer from "@/components/Footer";

export default function Home() {
  return (
    <main className="min-h-screen">
      <Hero />
      <CommandShowcase />
      <Features />
      <Architecture />
      <HowItWorks />
      <Footer />
    </main>
  );
}
