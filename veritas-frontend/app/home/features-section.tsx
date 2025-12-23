'use client';

import { motion } from 'framer-motion';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { fadeIn, stagger } from '@/lib/animations';

export function FeaturesSection() {
  return (
    <section className='py-20 bg-[#373b3b]/5'>
      <motion.div 
        initial='initial'
        whileInView='animate'
        viewport={{ once: true }}
        variants={stagger}
        className='container max-w-6xl mx-auto px-4'
      >
        <div className='mx-auto mb-16 max-w-3xl text-center'>
          <motion.h2 
            variants={fadeIn}
            className='text-3xl font-bold tracking-tight text-white sm:text-4xl lg:text-5xl'
          >
            Enterprise-Grade Risk Management
          </motion.h2>
          <motion.p 
            variants={fadeIn}
            className='mt-6 text-lg text-[#373b3b] max-w-2xl mx-auto'
          >
            Comprehensive suite of tools powered by advanced machine learning algorithms
          </motion.p>
        </div>
        <div className='mx-auto grid max-w-5xl gap-8 md:grid-cols-3'>
          {[
            {
              title: 'Real-Time Analysis',
              description: 'Continuous monitoring and assessment',
              features: ['Live risk scoring', 'Market condition tracking', 'Automated alerts']
            },
            {
              title: 'Scenario Testing',
              description: 'Advanced market simulations',
              features: ['Stress testing', 'Bull market scenarios', 'Risk forecasting']
            },
            {
              title: 'ML Intelligence',
              description: 'Smart predictive analytics',
              features: ['Pattern recognition', 'Anomaly detection', 'Trend analysis']
            }
          ].map((feature, index) => (
            <motion.div key={index} variants={fadeIn}>
              <Card className='border-2 border-[#373b3b]/20 bg-black'>
                <CardHeader>
                  <CardTitle className='text-white'>{feature.title}</CardTitle>
                  <CardDescription className='text-[#373b3b]'>
                    {feature.description}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <ul className='list-inside list-disc space-y-2 text-[#373b3b]'>
                    {feature.features.map((item, i) => (
                      <li key={i}>{item}</li>
                    ))}
                  </ul>
                </CardContent>
              </Card>
            </motion.div>
          ))}
        </div>
      </motion.div>
    </section>
  );
}
