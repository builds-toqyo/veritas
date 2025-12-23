"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Nav } from "@/components/nav";
import { fadeIn, stagger, scaleIn } from "@/lib/animations";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col bg-black text-white">
      <Nav />

      <main className="flex-1 overflow-hidden">
        <section className="relative min-h-[80vh] flex items-center justify-center py-20">
          <div className="absolute inset-0 bg-[url('/grid.svg')] opacity-10 bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]" />
          <motion.div 
            initial="initial"
            animate="animate"
            variants={stagger}
            className="container max-w-6xl mx-auto px-4"
          >
            <div className="flex flex-col items-center gap-8 text-center mx-auto max-w-3xl">
              <motion.h1 
                variants={fadeIn}
                className="text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl lg:text-7xl text-white">
                Real-Time Risk Assessment for{" "}
                <span className="bg-gradient-to-r from-white to-zinc-400 bg-clip-text text-transparent">
                  Veritas Vault
                </span>
              </motion.h1>
              <motion.p 
                variants={fadeIn}
                className="mx-auto max-w-[700px] text-lg text-zinc-400 md:text-xl">
                Advanced ML-powered risk analysis and monitoring system for DeFi operations.
              </motion.p>
              <div className="mt-6">
                <motion.div variants={scaleIn}>
                  <Button 
                    asChild 
                    size="lg" 
                    className="bg-[#373b3b] text-white hover:bg-[#565f5f] transition-colors duration-300"
                  >
                    <Link href="/dashboard">Launch Dashboard</Link>
                  </Button>
                </motion.div>
              </div>
            </div>
          </motion.div>
        </section>

        <section className="border-y border-[#181919] bg-black py-20">
          <div className="container max-w-6xl mx-auto px-4">
            <div className="mx-auto grid gap-16 md:grid-cols-3 max-w-4xl">
              <div className="text-center">
                <h3 className="text-4xl font-bold text-white">99.9%</h3>
                <p className="mt-2 text-zinc-400">Uptime</p>
              </div>
              <div className="text-center">
                <h3 className="text-4xl font-bold text-white">24/7</h3>
                <p className="mt-2 text-zinc-400">Monitoring</p>
              </div>
              <div className="text-center">
                <h3 className="text-4xl font-bold text-white">ML v2.0</h3>
                <p className="mt-2 text-zinc-400">Latest Model</p>
              </div>
            </div>
          </div>
        </section>

        <section className="py-20 bg-zinc-900">
          <motion.div 
            initial="initial"
            whileInView="animate"
            viewport={{ once: true }}
            variants={stagger}
            className="container max-w-6xl mx-auto px-4"
          >
            <div className="mx-auto mb-16 max-w-3xl text-center">
              <motion.h2 
                variants={fadeIn}
                className="text-3xl font-bold tracking-tight text-white sm:text-4xl lg:text-5xl">
                Enterprise-Grade Risk Management
              </motion.h2>
              <motion.p 
                variants={fadeIn}
                className="mt-6 text-lg text-zinc-400 max-w-2xl mx-auto">
                Comprehensive suite of tools powered by advanced machine learning algorithms
              </motion.p>
            </div>
            <div className="mx-auto grid max-w-5xl gap-8 md:grid-cols-3">
              <motion.div variants={fadeIn}>
                <Card className="border-2 border-zinc-800 bg-black">
                  <CardHeader>
                    <CardTitle className="text-white">Real-Time Analysis</CardTitle>
                    <CardDescription className="text-zinc-400">
                      Continuous monitoring and assessment
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <ul className="list-inside list-disc space-y-2 text-zinc-400">
                      <li>Live risk scoring</li>
                      <li>Market condition tracking</li>
                      <li>Automated alerts</li>
                    </ul>
                  </CardContent>
                </Card>
              </motion.div>

              <motion.div variants={fadeIn}>
                <Card className="border-2 border-zinc-800 bg-black">
                  <CardHeader>
                    <CardTitle className="text-white">Scenario Testing</CardTitle>
                    <CardDescription className="text-zinc-400">
                      Advanced market simulations
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <ul className="list-inside list-disc space-y-2 text-zinc-400">
                      <li>Stress testing</li>
                      <li>Bull market scenarios</li>
                      <li>Risk forecasting</li>
                    </ul>
                  </CardContent>
                </Card>
              </motion.div>

              <motion.div variants={fadeIn}>
                <Card className="border-2 border-zinc-800 bg-black">
                  <CardHeader>
                    <CardTitle className="text-white">ML Intelligence</CardTitle>
                    <CardDescription className="text-zinc-400">
                      Smart predictive analytics
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <ul className="list-inside list-disc space-y-2 text-zinc-400">
                      <li>Pattern recognition</li>
                      <li>Anomaly detection</li>
                      <li>Trend analysis</li>
                    </ul>
                  </CardContent>
                </Card>
              </motion.div>
            </div>
          </motion.div>
        </section>

        {/* Technology Stack Section */}
        <section className="py-20 bg-zinc-900">
          <div className="container max-w-6xl mx-auto px-4">
            <motion.div
              initial="initial"
              whileInView="animate"
              viewport={{ once: true }}
              variants={stagger}
              className="text-center mb-16"
            >
              <motion.h2 
                variants={fadeIn}
                className="text-3xl font-bold tracking-tight text-white sm:text-4xl lg:text-5xl mb-6"
              >
                Powered by Advanced Technology
              </motion.h2>
              <motion.p 
                variants={fadeIn}
                className="text-lg text-zinc-400 max-w-2xl mx-auto"
              >
                Built with cutting-edge technologies for optimal performance and reliability
              </motion.p>
            </motion.div>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
              <motion.div 
                variants={fadeIn}
                className="flex flex-col items-center p-6 bg-black border border-zinc-800 rounded-lg"
              >
                <div className="text-2xl font-bold text-white mb-2">ML Engine</div>
                <p className="text-zinc-400 text-center">Advanced predictive modeling</p>
              </motion.div>

              <motion.div 
                variants={fadeIn}
                className="flex flex-col items-center p-6 bg-black border border-zinc-800 rounded-lg"
              >
                <div className="text-2xl font-bold text-white mb-2">Smart Contracts</div>
                <p className="text-zinc-400 text-center">Secure and audited protocols</p>
              </motion.div>

              <motion.div 
                variants={fadeIn}
                className="flex flex-col items-center p-6 bg-black border border-zinc-800 rounded-lg"
              >
                <div className="text-2xl font-bold text-white mb-2">Real-time Data</div>
                <p className="text-zinc-400 text-center">Live market analysis</p>
              </motion.div>

              <motion.div 
                variants={fadeIn}
                className="flex flex-col items-center p-6 bg-black border border-zinc-800 rounded-lg"
              >
                <div className="text-2xl font-bold text-white mb-2">API Integration</div>
                <p className="text-zinc-400 text-center">Seamless connectivity</p>
              </motion.div>
            </div>
          </div>
        </section>

        {/* Security Section */}
        <section className="py-20 bg-black">
          <div className="container max-w-6xl mx-auto px-4">
            <div className="grid md:grid-cols-2 gap-12 items-center">
              <motion.div
                initial="initial"
                whileInView="animate"
                viewport={{ once: true }}
                variants={stagger}
                className="space-y-6"
              >
                <motion.h2 
                  variants={fadeIn}
                  className="text-3xl font-bold tracking-tight text-white sm:text-4xl"
                >
                  Enterprise-Grade Security
                </motion.h2>
                <motion.p 
                  variants={fadeIn}
                  className="text-lg text-zinc-400"
                >
                  Our platform implements multiple layers of security to protect your assets and data:
                </motion.p>
                <motion.ul 
                  variants={stagger}
                  className="space-y-4"
                >
                  {[
                    'Multi-signature authentication',
                    'Real-time threat monitoring',
                    'Automated risk assessment',
                    'Regular security audits'
                  ].map((item, index) => (
                    <motion.li
                      key={index}
                      variants={fadeIn}
                      className="flex items-center text-zinc-400"
                    >
                      <span className="w-2 h-2 bg-[#565f5f] rounded-full mr-3" />
                      {item}
                    </motion.li>
                  ))}
                </motion.ul>
              </motion.div>

              <motion.div
                initial="initial"
                whileInView="animate"
                viewport={{ once: true }}
                variants={fadeIn}
                className="grid grid-cols-2 gap-4"
              >
                {[
                  { title: '256-bit', desc: 'Encryption' },
                  { title: '24/7', desc: 'Monitoring' },
                  { title: '99.9%', desc: 'Uptime' },
                  { title: 'ISO 27001', desc: 'Compliant' }
                ].map((item, index) => (
                  <div 
                    key={index}
                    className="p-6 bg-zinc-900 border border-zinc-800 rounded-lg text-center"
                  >
                    <div className="text-2xl font-bold text-white mb-1">{item.title}</div>
                    <div className="text-zinc-400">{item.desc}</div>
                  </div>
                ))}
              </motion.div>
            </div>
          </div>
        </section>

        {/* How It Works Section */}
        <section className="py-20 bg-zinc-900">
          <div className="container max-w-6xl mx-auto px-4">
            <motion.div
              initial="initial"
              whileInView="animate"
              viewport={{ once: true }}
              variants={stagger}
              className="text-center mb-16"
            >
              <motion.h2
                variants={fadeIn}
                className="text-3xl font-bold tracking-tight text-white sm:text-4xl lg:text-5xl mb-6"
              >
                How Veritas Works
              </motion.h2>
              <motion.p
                variants={fadeIn}
                className="text-lg text-zinc-400 max-w-2xl mx-auto"
              >
                Our advanced risk management system operates through a seamless four-step process
              </motion.p>
            </motion.div>

            <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
              {[
                {
                  title: 'Data Collection',
                  description: 'Continuous monitoring of market conditions, liquidity pools, and risk factors',
                  number: '01'
                },
                {
                  title: 'ML Analysis',
                  description: 'Advanced algorithms process data to identify patterns and potential risks',
                  number: '02'
                },
                {
                  title: 'Risk Assessment',
                  description: 'Comprehensive evaluation of current market position and potential threats',
                  number: '03'
                },
                {
                  title: 'Smart Execution',
                  description: 'Automated risk mitigation strategies implemented through smart contracts',
                  number: '04'
                }
              ].map((step, index) => (
                <motion.div
                  key={index}
                  variants={fadeIn}
                  className="relative p-6 bg-black border border-zinc-800 rounded-lg"
                >
                  <div className="text-4xl font-bold text-zinc-600 mb-4">{step.number}</div>
                  <h3 className="text-xl font-bold text-white mb-2">{step.title}</h3>
                  <p className="text-zinc-400">{step.description}</p>
                  {index < 3 && (
                    <div className="hidden lg:block absolute -right-4 top-1/2 transform -translate-y-1/2 z-10">
                      <div className="w-8 h-[2px] bg-[#373b3b]"></div>
                    </div>
                  )}
                </motion.div>
              ))}
            </div>
          </div>
        </section>

        {/* Testimonials Section */}
        <section className="py-20 bg-black">
          <div className="container max-w-6xl mx-auto px-4">
            <motion.div
              initial="initial"
              whileInView="animate"
              viewport={{ once: true }}
              variants={stagger}
              className="text-center mb-16"
            >
              <motion.h2
                variants={fadeIn}
                className="text-3xl font-bold tracking-tight text-white sm:text-4xl lg:text-5xl mb-6"
              >
                Trusted by Industry Leaders
              </motion.h2>
              <motion.p
                variants={fadeIn}
                className="text-lg text-zinc-400 max-w-2xl mx-auto"
              >
                See what financial institutions and DeFi protocols say about Veritas
              </motion.p>
            </motion.div>

            <div className="grid md:grid-cols-3 gap-8">
              {[
                {
                  quote: 'Veritas has transformed how we approach risk management in DeFi. The ML-powered insights are invaluable.',
                  author: 'Sarah Chen',
                  role: 'Risk Manager, DeFi Capital',
                },
                {
                  quote: 'The real-time monitoring and automated risk assessment have helped us maintain optimal portfolio performance.',
                  author: 'Michael Rodriguez',
                  role: 'CTO, Blockchain Ventures',
                },
                {
                  quote: 'Implementing Veritas has significantly improved our risk-adjusted returns and operational efficiency.',
                  author: 'Emily Watson',
                  role: 'Head of Trading, Crypto Assets Ltd',
                },
              ].map((testimonial, index) => (
                <motion.div
                  key={index}
                  variants={fadeIn}
                  className="p-6 bg-zinc-900 border border-zinc-800 rounded-lg"
                >
                  <div className="mb-4 text-zinc-600">
                    <svg
                      className="w-8 h-8"
                      fill="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path d="M14.017 21v-7.391c0-5.704 3.731-9.57 8.983-10.609l.995 2.151c-2.432.917-3.995 3.638-3.995 5.849h4v10h-9.983zm-14.017 0v-7.391c0-5.704 3.748-9.57 9-10.609l.996 2.151c-2.433.917-3.996 3.638-3.996 5.849h3.983v10h-9.983z" />
                    </svg>
                  </div>
                  <p className="text-zinc-400 mb-4">{testimonial.quote}</p>
                  <div className="border-t border-zinc-800 pt-4">
                    <div className="font-bold text-white">{testimonial.author}</div>
                    <div className="text-sm text-zinc-400">{testimonial.role}</div>
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
        </section>

        {/* Partners Section */}
        <section className="py-20 bg-zinc-900">
          <div className="container max-w-6xl mx-auto px-4">
            <motion.div
              initial="initial"
              whileInView="animate"
              viewport={{ once: true }}
              variants={stagger}
              className="text-center mb-16"
            >
              <motion.h2
                variants={fadeIn}
                className="text-3xl font-bold tracking-tight text-white sm:text-4xl lg:text-5xl mb-6"
              >
                Our Partners
              </motion.h2>
              <motion.p
                variants={fadeIn}
                className="text-lg text-zinc-400 max-w-2xl mx-auto"
              >
                Working with leading organizations in blockchain and finance
              </motion.p>
            </motion.div>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
              {[
                'Ethereum Foundation',
                'Mantle Network',
                'Risk DAO',
                'DeFi Alliance',
              ].map((partner, index) => (
                <motion.div
                  key={index}
                  variants={fadeIn}
                  className="flex items-center justify-center p-6 bg-black border border-zinc-800 rounded-lg"
                >
                  <span className="text-lg font-bold text-white">{partner}</span>
                </motion.div>
              ))}
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="border-t border-[#181919] bg-black py-32 relative">
          <div className="absolute inset-0 bg-[url('/grid.svg')] opacity-5 bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]" />
          <div className="container max-w-6xl mx-auto px-4 flex flex-col items-center text-center">
            <motion.h2 
              initial="initial"
              whileInView="animate"
              viewport={{ once: true }}
              variants={fadeIn}
              className="text-3xl font-bold tracking-tight text-white sm:text-4xl">
              Ready to enhance your risk management?
            </motion.h2>
            <motion.p 
              initial="initial"
              whileInView="animate"
              viewport={{ once: true }}
              variants={fadeIn}
              className="mt-4 max-w-2xl text-lg text-zinc-400">
              Get started with Veritas today and access enterprise-grade risk assessment tools.
            </motion.p>
            <Button 
              asChild 
              size="lg" 
              className="mt-8 bg-[#373b3b] text-white hover:bg-[#565f5f] transition-colors duration-300"
            >
              <Link href="/dashboard">Access Dashboard</Link>
            </Button>
          </div>
        </section>
      </main>

      <footer className="border-t border-[#181919] py-8">
        <div className="container">
          <p className="text-center text-sm text-zinc-400">
            © 2025 Veritas. Built with advanced ML models for optimal risk management.
          </p>
        </div>
      </footer>
    </div>
  );
}
